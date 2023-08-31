package service

import (
	"context"
	"go.uber.org/zap"
	"time"
	"toktik-chat/internal/dao/mysql"
	"toktik-chat/internal/dao/redis"
	"toktik-chat/internal/model"
	"toktik-chat/internal/model/auto"
	"toktik-chat/internal/repo"
	"toktik-chat/pkg/myerr"
	"toktik-chat/pkg/rpc/client"
	"toktik-common/errcode"
	"toktik-rpc/kitex_gen/chat"
	"toktik-rpc/kitex_gen/interaction"
)

// ChatServiceImpl implements the last service interface defined in the IDL.
type ChatServiceImpl struct {
	respRepo    repo.HandlerResp
	cache       repo.Cache
	ChatRepo    repo.ChatRepo
	rClient     repo.RClientRepo
	transaction mysql.Transaction
}

func NewUserService() *ChatServiceImpl {
	return &ChatServiceImpl{
		cache:       redis.GetRdbCache(),
		respRepo:    repo.NewHandlerResps(),
		ChatRepo:    mysql.NewChatDao(),
		rClient:     redis.NewChatRClient(),
		transaction: mysql.NewTransaction(),
	}
}

// ChatAction implements the ChatServiceImpl interface.
func (cs *ChatServiceImpl) ChatAction(ctx context.Context, req *chat.ChatActionRequest) (resp *chat.ChatActionResponse, err error) {
	// 1.处理业务逻辑
	if req.Content == "" {
		return cs.respRepo.ChatActionResponse(myerr.MessageCanNotEmpty, model.MsgNil, &chat.ChatActionResponse{}), nil
	}
	// 判断target_id是否是自己的好友
	isFriendResp, err := client.InteractionClient.IsFriend(ctx, &interaction.IsFriendRequest{UserId: req.UserId, TargetId: req.ToUserId})
	if isFriendResp == nil {
		zap.L().Error("client.InteractionClient.IsFriend 返回空指针")
		return cs.respRepo.ChatActionResponse(errcode.ErrServer, model.MsgNil, &chat.ChatActionResponse{}), nil
	}
	if isFriendResp.StatusCode != model.RpcSuccess {
		zap.L().Error("client.InteractionClient.IsFriend err:", zap.Error(err))
		return cs.respRepo.ChatActionResponse(errcode.CreateErr(isFriendResp.StatusCode, model.MsgNil), isFriendResp.StatusMsg, &chat.ChatActionResponse{}), nil
	}
	if !isFriendResp.IsFriend {
		return cs.respRepo.ChatActionResponse(myerr.IsNotFriend, model.MsgNil, &chat.ChatActionResponse{}), nil
	}
	// 2.处理业务
	// 将消息存入 redis的list中，做两个缓存，一个缓存历史消息，一个缓存未持久化的数据库的消息
	// ①
	historyKey := auto.CreateChatHistoryKey(req.UserId, req.ToUserId)
	messageKey := auto.CreateChatMessageKey(req.UserId, req.ToUserId)
	//content := auto.CreateMessageContent(req.Content)
	timeNow := time.Now().Unix()
	content := auto.CreateMessageContent(req.UserId, req.ToUserId, req.Content)
	err = cs.rClient.PushHistoryMessage(ctx, historyKey, float64(timeNow), content)
	if err != nil {
		zap.L().Error("cs.rClient.PushHistoryMessage err:", zap.Error(err))
		return cs.respRepo.ChatActionResponse(errcode.ErrRedis, err.Error(), &chat.ChatActionResponse{}), nil
	}
	err = cs.rClient.PushDBMessage(ctx, messageKey, float64(timeNow), content)
	if err != nil {
		zap.L().Error("cs.rClient.PushDBMessage err:", zap.Error(err))
		return cs.respRepo.ChatActionResponse(errcode.ErrRedis, err.Error(), &chat.ChatActionResponse{}), nil
	}

	return cs.respRepo.ChatActionResponse(errcode.StatusOK, model.MsgNil, &chat.ChatActionResponse{}), nil
}

// MessageList implements the ChatServiceImpl interface.
func (cs *ChatServiceImpl) MessageList(ctx context.Context, req *chat.MessageListRequest) (resp *chat.MessageListResponse, err error) {
	resp = new(chat.MessageListResponse)
	// 1.处理业务逻辑
	// 2.处理业务
	// 从redis取消息，后发的消息先取
	historyKey := auto.CreateChatHistoryKey(req.UserId, req.ToUserId)
	messageList, err := cs.rClient.ZRangeMessageList(ctx, historyKey)
	if err != nil {
		zap.L().Error("cs.rClient.ZRangeMessageList err:", zap.Error(err))
		return cs.respRepo.MessageListResponse(errcode.ErrRedis, err.Error(), &chat.MessageListResponse{}), nil
	}
	if messageList == nil {
		// 未命中缓存，取数据库取消息
		messageList, err = cs.ChatRepo.GetMessageList(ctx, req.UserId, req.ToUserId)
		if err != nil {
			zap.L().Error("cs.rClient.GetMessageList err:", zap.Error(err))
			return cs.respRepo.MessageListResponse(errcode.ErrDB, err.Error(), &chat.MessageListResponse{}), nil
		}
		// 添加缓存
		if err = cs.rClient.PushManyHistoryMessage(ctx, historyKey, messageList); err != nil {
			zap.L().Error("cs.rClient.PushManyHistoryMessage err:", zap.Error(err))
			return cs.respRepo.MessageListResponse(errcode.ErrRedis, err.Error(), &chat.MessageListResponse{}), nil
		}
	}

	resp.MessageList = make([]*chat.Message, len(messageList))
	for i, message := range messageList {
		createdTime := message.CreatedAt.Unix()
		resp.MessageList[i] = &chat.Message{
			Id:         int64(i),
			ToUserId:   int64(message.ToUserId),
			FromUserId: int64(message.UserId),
			Content:    message.Content,
			CreateTime: createdTime,
		}
	}

	return cs.respRepo.MessageListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}
