// Code generated by Fastpb v0.0.2. DO NOT EDIT.

package interaction

import (
	fmt "fmt"
	fastpb "github.com/cloudwego/fastpb"
	user "toktik-rpc/kitex_gen/user"
)

var (
	_ = fmt.Errorf
	_ = fastpb.Skip
)

func (x *FollowActionRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 4:
		offset, err = x.fastReadField4(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FollowActionRequest[number], err)
}

func (x *FollowActionRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.Token, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *FollowActionRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.ToUserId, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *FollowActionRequest) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.ActionType, offset, err = fastpb.ReadInt32(buf, _type)
	return offset, err
}

func (x *FollowActionRequest) fastReadField4(buf []byte, _type int8) (offset int, err error) {
	x.MyUserId, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *FollowActionResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FollowActionResponse[number], err)
}

func (x *FollowActionResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.StatusCode, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *FollowActionResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.StatusMsg, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *FollowListRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FollowListRequest[number], err)
}

func (x *FollowListRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *FollowListRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.Token, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *FollowListResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FollowListResponse[number], err)
}

func (x *FollowListResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.StatusCode, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *FollowListResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.StatusMsg, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *FollowListResponse) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	var v user.User
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.UserList = append(x.UserList, &v)
	return offset, nil
}

func (x *FansListRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FansListRequest[number], err)
}

func (x *FansListRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *FansListRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.Token, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *FansListResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FansListResponse[number], err)
}

func (x *FansListResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.StatusCode, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *FansListResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.StatusMsg, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *FansListResponse) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	var v user.User
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.UserList = append(x.UserList, &v)
	return offset, nil
}

func (x *FriendListRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FriendListRequest[number], err)
}

func (x *FriendListRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *FriendListRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.Token, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *FriendListResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_FriendListResponse[number], err)
}

func (x *FriendListResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.StatusCode, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *FriendListResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.StatusMsg, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *FriendListResponse) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	var v user.User
	offset, err = fastpb.ReadMessage(buf, _type, &v)
	if err != nil {
		return offset, err
	}
	x.UserList = append(x.UserList, &v)
	return offset, nil
}

func (x *IsFollowTargetRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_IsFollowTargetRequest[number], err)
}

func (x *IsFollowTargetRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *IsFollowTargetRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.TargetId, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *IsFollowTargetResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_IsFollowTargetResponse[number], err)
}

func (x *IsFollowTargetResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.StatusCode, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *IsFollowTargetResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.StatusMsg, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *IsFollowTargetResponse) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.Exist, offset, err = fastpb.ReadBool(buf, _type)
	return offset, err
}

func (x *IsFollowManyTargetsRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_IsFollowManyTargetsRequest[number], err)
}

func (x *IsFollowManyTargetsRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *IsFollowManyTargetsRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	offset, err = fastpb.ReadList(buf, _type,
		func(buf []byte, _type int8) (n int, err error) {
			var v int64
			v, offset, err = fastpb.ReadInt64(buf, _type)
			if err != nil {
				return offset, err
			}
			x.TargetIds = append(x.TargetIds, v)
			return offset, err
		})
	return offset, err
}

func (x *IsFollowManyTargetsResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_IsFollowManyTargetsResponse[number], err)
}

func (x *IsFollowManyTargetsResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.StatusCode, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *IsFollowManyTargetsResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.StatusMsg, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *IsFollowManyTargetsResponse) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	offset, err = fastpb.ReadList(buf, _type,
		func(buf []byte, _type int8) (n int, err error) {
			var v bool
			v, offset, err = fastpb.ReadBool(buf, _type)
			if err != nil {
				return offset, err
			}
			x.ManyExist = append(x.ManyExist, v)
			return offset, err
		})
	return offset, err
}

func (x *IsFriendRequest) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_IsFriendRequest[number], err)
}

func (x *IsFriendRequest) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *IsFriendRequest) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.TargetId, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *IsFriendResponse) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 3:
		offset, err = x.fastReadField3(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_IsFriendResponse[number], err)
}

func (x *IsFriendResponse) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.StatusCode, offset, err = fastpb.ReadInt64(buf, _type)
	return offset, err
}

func (x *IsFriendResponse) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.StatusMsg, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *IsFriendResponse) fastReadField3(buf []byte, _type int8) (offset int, err error) {
	x.IsFriend, offset, err = fastpb.ReadBool(buf, _type)
	return offset, err
}

func (x *FollowActionRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	offset += x.fastWriteField4(buf[offset:])
	return offset
}

func (x *FollowActionRequest) fastWriteField1(buf []byte) (offset int) {
	if x.Token == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 1, x.GetToken())
	return offset
}

func (x *FollowActionRequest) fastWriteField2(buf []byte) (offset int) {
	if x.ToUserId == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 2, x.GetToUserId())
	return offset
}

func (x *FollowActionRequest) fastWriteField3(buf []byte) (offset int) {
	if x.ActionType == 0 {
		return offset
	}
	offset += fastpb.WriteInt32(buf[offset:], 3, x.GetActionType())
	return offset
}

func (x *FollowActionRequest) fastWriteField4(buf []byte) (offset int) {
	if x.MyUserId == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 4, x.GetMyUserId())
	return offset
}

func (x *FollowActionResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *FollowActionResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetStatusCode())
	return offset
}

func (x *FollowActionResponse) fastWriteField2(buf []byte) (offset int) {
	if x.StatusMsg == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetStatusMsg())
	return offset
}

func (x *FollowListRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *FollowListRequest) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *FollowListRequest) fastWriteField2(buf []byte) (offset int) {
	if x.Token == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetToken())
	return offset
}

func (x *FollowListResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *FollowListResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetStatusCode())
	return offset
}

func (x *FollowListResponse) fastWriteField2(buf []byte) (offset int) {
	if x.StatusMsg == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetStatusMsg())
	return offset
}

func (x *FollowListResponse) fastWriteField3(buf []byte) (offset int) {
	if x.UserList == nil {
		return offset
	}
	for i := range x.GetUserList() {
		offset += fastpb.WriteMessage(buf[offset:], 3, x.GetUserList()[i])
	}
	return offset
}

func (x *FansListRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *FansListRequest) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *FansListRequest) fastWriteField2(buf []byte) (offset int) {
	if x.Token == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetToken())
	return offset
}

func (x *FansListResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *FansListResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetStatusCode())
	return offset
}

func (x *FansListResponse) fastWriteField2(buf []byte) (offset int) {
	if x.StatusMsg == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetStatusMsg())
	return offset
}

func (x *FansListResponse) fastWriteField3(buf []byte) (offset int) {
	if x.UserList == nil {
		return offset
	}
	for i := range x.GetUserList() {
		offset += fastpb.WriteMessage(buf[offset:], 3, x.GetUserList()[i])
	}
	return offset
}

func (x *FriendListRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *FriendListRequest) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *FriendListRequest) fastWriteField2(buf []byte) (offset int) {
	if x.Token == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetToken())
	return offset
}

func (x *FriendListResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *FriendListResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetStatusCode())
	return offset
}

func (x *FriendListResponse) fastWriteField2(buf []byte) (offset int) {
	if x.StatusMsg == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetStatusMsg())
	return offset
}

func (x *FriendListResponse) fastWriteField3(buf []byte) (offset int) {
	if x.UserList == nil {
		return offset
	}
	for i := range x.GetUserList() {
		offset += fastpb.WriteMessage(buf[offset:], 3, x.GetUserList()[i])
	}
	return offset
}

func (x *IsFollowTargetRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *IsFollowTargetRequest) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *IsFollowTargetRequest) fastWriteField2(buf []byte) (offset int) {
	if x.TargetId == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 2, x.GetTargetId())
	return offset
}

func (x *IsFollowTargetResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *IsFollowTargetResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetStatusCode())
	return offset
}

func (x *IsFollowTargetResponse) fastWriteField2(buf []byte) (offset int) {
	if x.StatusMsg == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetStatusMsg())
	return offset
}

func (x *IsFollowTargetResponse) fastWriteField3(buf []byte) (offset int) {
	if !x.Exist {
		return offset
	}
	offset += fastpb.WriteBool(buf[offset:], 3, x.GetExist())
	return offset
}

func (x *IsFollowManyTargetsRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *IsFollowManyTargetsRequest) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *IsFollowManyTargetsRequest) fastWriteField2(buf []byte) (offset int) {
	if len(x.TargetIds) == 0 {
		return offset
	}
	offset += fastpb.WriteListPacked(buf[offset:], 2, len(x.GetTargetIds()),
		func(buf []byte, numTagOrKey, numIdxOrVal int32) int {
			offset := 0
			offset += fastpb.WriteInt64(buf[offset:], numTagOrKey, x.GetTargetIds()[numIdxOrVal])
			return offset
		})
	return offset
}

func (x *IsFollowManyTargetsResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *IsFollowManyTargetsResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetStatusCode())
	return offset
}

func (x *IsFollowManyTargetsResponse) fastWriteField2(buf []byte) (offset int) {
	if x.StatusMsg == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetStatusMsg())
	return offset
}

func (x *IsFollowManyTargetsResponse) fastWriteField3(buf []byte) (offset int) {
	if len(x.ManyExist) == 0 {
		return offset
	}
	offset += fastpb.WriteListPacked(buf[offset:], 3, len(x.GetManyExist()),
		func(buf []byte, numTagOrKey, numIdxOrVal int32) int {
			offset := 0
			offset += fastpb.WriteBool(buf[offset:], numTagOrKey, x.GetManyExist()[numIdxOrVal])
			return offset
		})
	return offset
}

func (x *IsFriendRequest) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *IsFriendRequest) fastWriteField1(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetUserId())
	return offset
}

func (x *IsFriendRequest) fastWriteField2(buf []byte) (offset int) {
	if x.TargetId == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 2, x.GetTargetId())
	return offset
}

func (x *IsFriendResponse) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	offset += x.fastWriteField3(buf[offset:])
	return offset
}

func (x *IsFriendResponse) fastWriteField1(buf []byte) (offset int) {
	if x.StatusCode == 0 {
		return offset
	}
	offset += fastpb.WriteInt64(buf[offset:], 1, x.GetStatusCode())
	return offset
}

func (x *IsFriendResponse) fastWriteField2(buf []byte) (offset int) {
	if x.StatusMsg == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 2, x.GetStatusMsg())
	return offset
}

func (x *IsFriendResponse) fastWriteField3(buf []byte) (offset int) {
	if !x.IsFriend {
		return offset
	}
	offset += fastpb.WriteBool(buf[offset:], 3, x.GetIsFriend())
	return offset
}

func (x *FollowActionRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	n += x.sizeField4()
	return n
}

func (x *FollowActionRequest) sizeField1() (n int) {
	if x.Token == "" {
		return n
	}
	n += fastpb.SizeString(1, x.GetToken())
	return n
}

func (x *FollowActionRequest) sizeField2() (n int) {
	if x.ToUserId == 0 {
		return n
	}
	n += fastpb.SizeInt64(2, x.GetToUserId())
	return n
}

func (x *FollowActionRequest) sizeField3() (n int) {
	if x.ActionType == 0 {
		return n
	}
	n += fastpb.SizeInt32(3, x.GetActionType())
	return n
}

func (x *FollowActionRequest) sizeField4() (n int) {
	if x.MyUserId == 0 {
		return n
	}
	n += fastpb.SizeInt64(4, x.GetMyUserId())
	return n
}

func (x *FollowActionResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *FollowActionResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetStatusCode())
	return n
}

func (x *FollowActionResponse) sizeField2() (n int) {
	if x.StatusMsg == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetStatusMsg())
	return n
}

func (x *FollowListRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *FollowListRequest) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetUserId())
	return n
}

func (x *FollowListRequest) sizeField2() (n int) {
	if x.Token == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetToken())
	return n
}

func (x *FollowListResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *FollowListResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetStatusCode())
	return n
}

func (x *FollowListResponse) sizeField2() (n int) {
	if x.StatusMsg == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetStatusMsg())
	return n
}

func (x *FollowListResponse) sizeField3() (n int) {
	if x.UserList == nil {
		return n
	}
	for i := range x.GetUserList() {
		n += fastpb.SizeMessage(3, x.GetUserList()[i])
	}
	return n
}

func (x *FansListRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *FansListRequest) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetUserId())
	return n
}

func (x *FansListRequest) sizeField2() (n int) {
	if x.Token == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetToken())
	return n
}

func (x *FansListResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *FansListResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetStatusCode())
	return n
}

func (x *FansListResponse) sizeField2() (n int) {
	if x.StatusMsg == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetStatusMsg())
	return n
}

func (x *FansListResponse) sizeField3() (n int) {
	if x.UserList == nil {
		return n
	}
	for i := range x.GetUserList() {
		n += fastpb.SizeMessage(3, x.GetUserList()[i])
	}
	return n
}

func (x *FriendListRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *FriendListRequest) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetUserId())
	return n
}

func (x *FriendListRequest) sizeField2() (n int) {
	if x.Token == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetToken())
	return n
}

func (x *FriendListResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *FriendListResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetStatusCode())
	return n
}

func (x *FriendListResponse) sizeField2() (n int) {
	if x.StatusMsg == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetStatusMsg())
	return n
}

func (x *FriendListResponse) sizeField3() (n int) {
	if x.UserList == nil {
		return n
	}
	for i := range x.GetUserList() {
		n += fastpb.SizeMessage(3, x.GetUserList()[i])
	}
	return n
}

func (x *IsFollowTargetRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *IsFollowTargetRequest) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetUserId())
	return n
}

func (x *IsFollowTargetRequest) sizeField2() (n int) {
	if x.TargetId == 0 {
		return n
	}
	n += fastpb.SizeInt64(2, x.GetTargetId())
	return n
}

func (x *IsFollowTargetResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *IsFollowTargetResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetStatusCode())
	return n
}

func (x *IsFollowTargetResponse) sizeField2() (n int) {
	if x.StatusMsg == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetStatusMsg())
	return n
}

func (x *IsFollowTargetResponse) sizeField3() (n int) {
	if !x.Exist {
		return n
	}
	n += fastpb.SizeBool(3, x.GetExist())
	return n
}

func (x *IsFollowManyTargetsRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *IsFollowManyTargetsRequest) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetUserId())
	return n
}

func (x *IsFollowManyTargetsRequest) sizeField2() (n int) {
	if len(x.TargetIds) == 0 {
		return n
	}
	n += fastpb.SizeListPacked(2, len(x.GetTargetIds()),
		func(numTagOrKey, numIdxOrVal int32) int {
			n := 0
			n += fastpb.SizeInt64(numTagOrKey, x.GetTargetIds()[numIdxOrVal])
			return n
		})
	return n
}

func (x *IsFollowManyTargetsResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *IsFollowManyTargetsResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetStatusCode())
	return n
}

func (x *IsFollowManyTargetsResponse) sizeField2() (n int) {
	if x.StatusMsg == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetStatusMsg())
	return n
}

func (x *IsFollowManyTargetsResponse) sizeField3() (n int) {
	if len(x.ManyExist) == 0 {
		return n
	}
	n += fastpb.SizeListPacked(3, len(x.GetManyExist()),
		func(numTagOrKey, numIdxOrVal int32) int {
			n := 0
			n += fastpb.SizeBool(numTagOrKey, x.GetManyExist()[numIdxOrVal])
			return n
		})
	return n
}

func (x *IsFriendRequest) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *IsFriendRequest) sizeField1() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetUserId())
	return n
}

func (x *IsFriendRequest) sizeField2() (n int) {
	if x.TargetId == 0 {
		return n
	}
	n += fastpb.SizeInt64(2, x.GetTargetId())
	return n
}

func (x *IsFriendResponse) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	n += x.sizeField3()
	return n
}

func (x *IsFriendResponse) sizeField1() (n int) {
	if x.StatusCode == 0 {
		return n
	}
	n += fastpb.SizeInt64(1, x.GetStatusCode())
	return n
}

func (x *IsFriendResponse) sizeField2() (n int) {
	if x.StatusMsg == "" {
		return n
	}
	n += fastpb.SizeString(2, x.GetStatusMsg())
	return n
}

func (x *IsFriendResponse) sizeField3() (n int) {
	if !x.IsFriend {
		return n
	}
	n += fastpb.SizeBool(3, x.GetIsFriend())
	return n
}

var fieldIDToName_FollowActionRequest = map[int32]string{
	1: "Token",
	2: "ToUserId",
	3: "ActionType",
	4: "MyUserId",
}

var fieldIDToName_FollowActionResponse = map[int32]string{
	1: "StatusCode",
	2: "StatusMsg",
}

var fieldIDToName_FollowListRequest = map[int32]string{
	1: "UserId",
	2: "Token",
}

var fieldIDToName_FollowListResponse = map[int32]string{
	1: "StatusCode",
	2: "StatusMsg",
	3: "UserList",
}

var fieldIDToName_FansListRequest = map[int32]string{
	1: "UserId",
	2: "Token",
}

var fieldIDToName_FansListResponse = map[int32]string{
	1: "StatusCode",
	2: "StatusMsg",
	3: "UserList",
}

var fieldIDToName_FriendListRequest = map[int32]string{
	1: "UserId",
	2: "Token",
}

var fieldIDToName_FriendListResponse = map[int32]string{
	1: "StatusCode",
	2: "StatusMsg",
	3: "UserList",
}

var fieldIDToName_IsFollowTargetRequest = map[int32]string{
	1: "UserId",
	2: "TargetId",
}

var fieldIDToName_IsFollowTargetResponse = map[int32]string{
	1: "StatusCode",
	2: "StatusMsg",
	3: "Exist",
}

var fieldIDToName_IsFollowManyTargetsRequest = map[int32]string{
	1: "UserId",
	2: "TargetIds",
}

var fieldIDToName_IsFollowManyTargetsResponse = map[int32]string{
	1: "StatusCode",
	2: "StatusMsg",
	3: "ManyExist",
}

var fieldIDToName_IsFriendRequest = map[int32]string{
	1: "UserId",
	2: "TargetId",
}

var fieldIDToName_IsFriendResponse = map[int32]string{
	1: "StatusCode",
	2: "StatusMsg",
	3: "IsFriend",
}

var _ = user.File_user_proto
