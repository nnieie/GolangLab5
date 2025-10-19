package main

import (
	"context"

	"github.com/nnieie/golanglab5/cmd/chat/service"
	chat "github.com/nnieie/golanglab5/kitex_gen/chat"
	"github.com/nnieie/golanglab5/pkg/utils"
)

// ChatServiceImpl implements the last service interface defined in the IDL.
type ChatServiceImpl struct{}

// SendPrivateMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) SendPrivateMessage(ctx context.Context, req *chat.SendPrivateMessageRequest,
) (resp *chat.SendPrivateMessageResponse, err error) {
	resp = new(chat.SendPrivateMessageResponse)
	err = service.NewChatService(ctx).SendPrivateMessage(req.Data.FromUserId, req.Data.ToUserId, req.Data.Content)
	resp.Base = utils.BuildBaseResp(err)
	return
}

// QueryPrivateOfflineMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) QueryPrivateOfflineMessage(ctx context.Context, req *chat.QueryPrivateOfflineMessageRequest,
) (resp *chat.QueryPrivateOfflineMessageResponse, err error) {
	resp = new(chat.QueryPrivateOfflineMessageResponse)
	msgs, err := service.NewChatService(ctx).GetOfflinePrivateMessage(req.UserId, req.ToUserId, req.PageNum, req.PageSize)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = msgs
	return
}

// QueryPrivateHistoryMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) QueryPrivateHistoryMessage(ctx context.Context, req *chat.QueryPrivateHistoryMessageRequest,
) (resp *chat.QueryPrivateHistoryMessageResponse, err error) {
	resp = new(chat.QueryPrivateHistoryMessageResponse)
	msgs, err := service.NewChatService(ctx).GetPrivateHistoryMessage(req.UserId, req.ToUserId, req.PageNum, req.PageSize)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = msgs
	return
}

// SendGroupMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) SendGroupMessage(ctx context.Context, req *chat.SendGroupMessageRequest,
) (resp *chat.SendGroupMessageResponse, err error) {
	resp = new(chat.SendGroupMessageResponse)
	err = service.NewChatService(ctx).SendGroupMessage(req.Data.FromUserId, req.Data.GroupId, req.Data.Content)
	resp.Base = utils.BuildBaseResp(err)
	return
}

// QueryGroupOfflineMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) QueryGroupOfflineMessage(ctx context.Context, req *chat.QueryGroupOfflineMessageRequest,
) (resp *chat.QueryGroupOfflineMessageResponse, err error) {
	resp = new(chat.QueryGroupOfflineMessageResponse)
	msgs, err := service.NewChatService(ctx).GetOfflineGroupMessage(req.UserId, req.GroupId, req.PageNum, req.PageSize)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = msgs
	return
}

// QueryGroupHistoryMessage implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) QueryGroupHistoryMessage(ctx context.Context, req *chat.QueryGroupHistoryMessageRequest,
) (resp *chat.QueryGroupHistoryMessageResponse, err error) {
	resp = new(chat.QueryGroupHistoryMessageResponse)
	msgs, err := service.NewChatService(ctx).GetGroupHistoryMessage(req.GroupId, req.PageNum, req.PageSize)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = msgs
	return
}

// QueryGroupMembers implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) QueryGroupMembers(ctx context.Context, req *chat.QueryGroupMembersRequest,
) (resp *chat.QueryGroupMembersResponse, err error) {
	resp = new(chat.QueryGroupMembersResponse)
	members, err := service.NewChatService(ctx).GetGroupMembers(req.GroupId)
	resp.Base = utils.BuildBaseResp(err)
	resp.Members = members
	return
}

// CheckUserExistInGroup implements the ChatServiceImpl interface.
func (s *ChatServiceImpl) CheckUserExistInGroup(ctx context.Context, req *chat.CheckUserExistInGroupRequest,
) (resp *chat.CheckUserExistInGroupResponse, err error) {
	resp = new(chat.CheckUserExistInGroupResponse)
	exist, err := service.NewChatService(ctx).CheckUserExistInGroup(req.UserId, req.GroupId)
	resp.Base = utils.BuildBaseResp(err)
	resp.Exist = &exist
	return
}
