package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"time"
	"toktik-chat/internal/cache"
	"toktik-chat/internal/dao/mysql"
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
	ChatRepo    repo.ChatRepo
	rCache      repo.RCacheRepo
	transaction mysql.Transaction
}

func NewChatService() *ChatServiceImpl {
	return &ChatServiceImpl{
		respRepo:    repo.NewHandlerResps(),
		ChatRepo:    mysql.NewChatDao(),
		rCache:      cache.NewChatrCache(),
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
	err = cs.rCache.PushHistoryMessage(ctx, historyKey, float64(timeNow), content)
	if err != nil {
		zap.L().Error("cs.rCache.PushHistoryMessage err:", zap.Error(err))
		return cs.respRepo.ChatActionResponse(errcode.ErrRedis, err.Error(), &chat.ChatActionResponse{}), nil
	}
	err = cs.rCache.PushDBMessage(ctx, messageKey, float64(timeNow), content)
	if err != nil {
		zap.L().Error("cs.rCache.PushDBMessage err:", zap.Error(err))
		return cs.respRepo.ChatActionResponse(errcode.ErrRedis, err.Error(), &chat.ChatActionResponse{}), nil
	}

	return cs.respRepo.ChatActionResponse(errcode.StatusOK, model.MsgNil, &chat.ChatActionResponse{}), nil
}

// MessageList implements the ChatServiceImpl interface.
func (cs *ChatServiceImpl) MessageList(ctx context.Context, req *chat.MessageListRequest) (resp *chat.MessageListResponse, err error) {
	resp = new(chat.MessageListResponse)
	req.PreMsgTime = req.PreMsgTime + 1
	// 1.处理业务逻辑
	// 2.处理业务
	// 从redis取消息，后发的消息先取
	historyKey := auto.CreateChatHistoryKey(req.UserId, req.ToUserId)
	messageList, err := cs.rCache.ZRangeMessageList(ctx, historyKey, req.PreMsgTime)
	if err != nil {
		zap.L().Error("cs.rCache.ZRangeMessageList err:", zap.Error(err))
		return cs.respRepo.MessageListResponse(errcode.ErrRedis, err.Error(), &chat.MessageListResponse{}), nil
	}
	exist, err := cs.rCache.KeyExist(ctx, historyKey)
	if err != nil {
		zap.L().Error("cs.rCache.KeyExist err:", zap.Error(err))
		return cs.respRepo.MessageListResponse(errcode.ErrRedis, err.Error(), &chat.MessageListResponse{}), nil
	}
	if !exist {
		fmt.Println("缓存没有，！！！！！！！！取 数据库！！！！！！！！！！！！！！！！！！！！！")
		// 未命中缓存，取数据库取消息
		messageList, err = cs.ChatRepo.GetMessageList(ctx, req.UserId, req.ToUserId, req.PreMsgTime)
		if err != nil {
			zap.L().Error("cs.rCache.GetMessageList err:", zap.Error(err))
			return cs.respRepo.MessageListResponse(errcode.ErrDB, err.Error(), &chat.MessageListResponse{}), nil
		}
		if messageList == nil {
			return cs.respRepo.MessageListResponse(errcode.StatusOK, model.MsgNil, &chat.MessageListResponse{}), nil
		}
		// 添加缓存
		if err = cs.rCache.PushManyHistoryMessage(ctx, historyKey, messageList); err != nil {
			zap.L().Error("cs.rCache.PushManyHistoryMessage err:", zap.Error(err))
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
			CreateTime: createdTime, // 将ms级转为 s级
		}
	}

	return cs.respRepo.MessageListResponse(errcode.StatusOK, model.MsgNil, resp), nil
}

func (cs *ChatServiceImpl) GetFriendLatestMessage(ctx context.Context, req *chat.GetFriendLatestMessageRequest) (resp *chat.GetFriendLatestMessageResponse, err error) {
	fmt.Printf("GetFriendLatestMessage req:%#v\n", req)
	resp = new(chat.GetFriendLatestMessageResponse)
	resp.MessageList = make([]string, len(req.FriendIds))
	resp.MsgTypeList = make([]int32, len(req.FriendIds))
	// 先去缓存查 message
	for i, friendId := range req.FriendIds {
		key := auto.CreateChatHistoryKey(req.UserId, friendId)
		message, msgType, err := cs.rCache.ZGetFriendLatestMessage(ctx, key, req.UserId)
		if err != nil {
			zap.L().Error("cs.rCache.ZGetFriendLatestMessage err:", zap.Error(err))
			return cs.respRepo.GetFriendLatestMessageResponse(errcode.ErrRedis, err.Error(), &chat.GetFriendLatestMessageResponse{}), nil
		}
		resp.MessageList[i] = message
		resp.MsgTypeList[i] = msgType
	}
	return cs.respRepo.GetFriendLatestMessageResponse(errcode.StatusOK, model.MsgNil, resp), nil
}
