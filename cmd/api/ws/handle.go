package chat

import (
	"encoding/json"

	"github.com/nnieie/golanglab5/cmd/api/biz/model/common"
	"github.com/nnieie/golanglab5/cmd/api/pack"
	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/utils"
)

const (
	TypeSendPrivateMessage    = 1
	TypePrivateHistoryMessage = 2
	TypePrivateOfflineMessage = 3
	TypeSendGroupMessage      = 4
	TypeGroupHistoryMessage   = 5
	TypeGroupOfflineMessage   = 6
)

// handleMessage 根据消息类型分发到对应的处理函数
func (c *Client) handleMessage(msg *common.Message) error {
	switch msg.Type {
	case TypeSendPrivateMessage:
		c.handlePrivateSendMsg(msg)
	case TypePrivateHistoryMessage:
		c.handlePrivateHistoryMsg(msg)
	case TypePrivateOfflineMessage:
		c.handlePrivateOfflineMsg(msg)
	case TypeSendGroupMessage:
		c.handleGroupSendMsg(msg)
	case TypeGroupHistoryMessage:
		c.handleGroupHistoryMsg(msg)
	case TypeGroupOfflineMessage:
		c.handleGroupOfflineMsg(msg)
	default:
		return errno.ParamErr
	}

	return nil
}

func (c *Client) handlePrivateSendMsg(msg *common.Message) {
	resp := new(common.SendPrivateMessageResponse)
	err := c.svc.SendPrivateMessage(msg.ToUserID, msg.Content)
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.sendStruct(resp)
		return
	}

	// 推送给接收者
	msg.FromUserID = c.userID
	payload, err := json.Marshal(msg)
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.sendStruct(resp)
		return
	}
	hub.broadcast <- &Broadcast{
		TargetUserIDs: []int64{msg.ToUserID},
		Payload:       payload,
	}
	resp.Base = pack.BuildBaseResp(nil)
	c.sendStruct(resp)
}

func (c *Client) handlePrivateHistoryMsg(msg *common.Message) {
	resp := new(common.QueryPrivateHistoryMessageResponse)
	msgs, err := c.svc.GetPrivateHistoryMessage(msg.ToUserID, msg.PageNum, msg.PageSize)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		c.sendStruct(resp)
		return
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.Messages = pack.PrivateMessagesRPCToPrivateMessages(msgs)
	c.sendStruct(resp)
}

func (c *Client) handlePrivateOfflineMsg(msg *common.Message) {
	resp := new(common.QueryPrivateOfflineMessageResponse)
	msgs, err := c.svc.GetPrivateOfflineMessage(msg.ToUserID, msg.PageNum, msg.PageSize)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		c.sendStruct(resp)
		return
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.Messages = pack.PrivateMessagesRPCToPrivateMessages(msgs)
	c.sendStruct(resp)
}

func (c *Client) handleGroupSendMsg(msg *common.Message) {
	resp := new(common.SendGroupMessageResponse)
	err := c.svc.SendGroupMessage(msg.GroupID, msg.Content)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		c.sendStruct(resp)
		return
	}

	// 通过 RPC 获取群成员列表
	ids, err := c.svc.GetGroupMembers(msg.GroupID)
	logger.Debugf("Group %d members: %v", msg.GroupID, ids)
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.sendStruct(resp)
		return
	}

	// 推送给接收者
	msg.FromUserID = c.userID
	payload, err := json.Marshal(msg)
	if err != nil {
		resp := utils.BuildBaseResp(err)
		c.sendStruct(resp)
		return
	}
	hub.broadcast <- &Broadcast{
		TargetUserIDs: ids,
		Payload:       payload,
	}
	resp.Base = pack.BuildBaseResp(nil)
	c.sendStruct(resp)
}

func (c *Client) handleGroupHistoryMsg(msg *common.Message) {
	resp := new(common.QueryGroupHistoryMessageResponse)
	msgs, err := c.svc.GetGroupHistoryMessage(msg.GroupID, msg.PageNum, msg.PageSize)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		c.sendStruct(resp)
		return
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.Messages = pack.GroupMessagesRPCToGroupMessages(msgs)
	c.sendStruct(resp)
}

func (c *Client) handleGroupOfflineMsg(msg *common.Message) {
	resp := new(common.QueryGroupOfflineMessageResponse)
	msgs, err := c.svc.GetGroupOfflineMessage(msg.GroupID, msg.PageNum, msg.PageSize)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		c.sendStruct(resp)
		return
	}

	resp.Base = pack.BuildBaseResp(nil)
	resp.Messages = pack.GroupMessagesRPCToGroupMessages(msgs)
	c.sendStruct(resp)
}
