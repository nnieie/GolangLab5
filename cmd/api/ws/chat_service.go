package chat

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/nnieie/golanglab5/cmd/api/biz/handler/mw/jwt"
	"github.com/nnieie/golanglab5/cmd/api/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
	"github.com/nnieie/golanglab5/kitex_gen/chat"
)

type ChatService struct {
	ctx    context.Context
	c      *app.RequestContext
	userID int64
}

func NewChatService(ctx context.Context, c *app.RequestContext) *ChatService {
	userID, err := jwt.ExtractUserID(c)
	if err != nil {
		return nil
	}
	return &ChatService{ctx: ctx, c: c, userID: userID}
}

func (s *ChatService) SendPrivateMessage(toUserID int64, content string) error {
	req := &chat.SendPrivateMessageRequest{
		Data: &base.PrivateMessage{
			FromUserId: s.userID,
			ToUserId:   toUserID,
			Content:    content,
		},
	}

	_, err := rpc.SendPrivateMessage(s.ctx, req)
	return err
}

func (s *ChatService) GetPrivateOfflineMessage(toUserID int64, pageNum, pageSize int64) ([]*base.PrivateMessage, error) {
	req := &chat.QueryPrivateOfflineMessageRequest{
		UserId:   s.userID,
		ToUserId: toUserID,
		PageNum:  pageNum,
		PageSize: pageSize,
	}

	resp, err := rpc.QueryPrivateOfflineMessages(s.ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (s *ChatService) GetPrivateHistoryMessage(toUserID int64, pageNum, pageSize int64) ([]*base.PrivateMessage, error) {
	req := &chat.QueryPrivateHistoryMessageRequest{
		UserId:   s.userID,
		ToUserId: toUserID,
		PageNum:  pageNum,
		PageSize: pageSize,
	}

	resp, err := rpc.QueryPrivateMessageHistory(s.ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (s *ChatService) SendGroupMessage(groupID int64, content string) error {
	req := &chat.SendGroupMessageRequest{
		Data: &base.GroupMessage{
			FromUserId: s.userID,
			GroupId:    groupID,
			Content:    content,
		},
	}

	_, err := rpc.SendGroupMessage(s.ctx, req)
	return err
}

func (s *ChatService) GetGroupOfflineMessage(groupID int64, pageNum, pageSize int64) ([]*base.GroupMessage, error) {
	req := &chat.QueryGroupOfflineMessageRequest{
		UserId:   s.userID,
		GroupId:  groupID,
		PageNum:  pageNum,
		PageSize: pageSize,
	}

	resp, err := rpc.QueryGroupOfflineMessages(s.ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (s *ChatService) GetGroupHistoryMessage(groupID int64, pageNum, pageSize int64) ([]*base.GroupMessage, error) {
	req := &chat.QueryGroupHistoryMessageRequest{
		UserId:   s.userID,
		GroupId:  groupID,
		PageNum:  pageNum,
		PageSize: pageSize,
	}

	resp, err := rpc.QueryGroupMessageHistory(s.ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (s *ChatService) GetGroupMembers(groupID int64) ([]int64, error) {
	req := &chat.QueryGroupMembersRequest{
		GroupId: groupID,
	}

	resp, err := rpc.QueryGroupMembers(s.ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Members, nil
}

func (s *ChatService) CheckUserExistInGroup(userID, groupID int64) (bool, error) {
	req := &chat.CheckUserExistInGroupRequest{
		UserId:  userID,
		GroupId: groupID,
	}

	resp, err := rpc.CheckUserExistInGroup(s.ctx, req)
	if err != nil {
		return false, err
	}

	return *resp.Exist, nil
}
