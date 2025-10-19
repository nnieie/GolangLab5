package pack

import (
	"github.com/nnieie/golanglab5/cmd/chat/dal/db"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func DBPrivateMessageToChatPrivateMessage(dbMsg *db.PrivateMessage) *base.PrivateMessage {
	if dbMsg == nil {
		return nil
	}
	return &base.PrivateMessage{
		FromUserId: dbMsg.FromUserID,
		ToUserId:   dbMsg.ToUserID,
		Content:    dbMsg.Content,
		CreatedAt:  dbMsg.CreatedAt.Unix(),
	}
}

func DBPrivateMessagesToChatPrivateMessages(dbMsgs []*db.PrivateMessage) []*base.PrivateMessage {
	chatMsgs := make([]*base.PrivateMessage, 0, len(dbMsgs))
	for _, msg := range dbMsgs {
		chatMsgs = append(chatMsgs, DBPrivateMessageToChatPrivateMessage(msg))
	}
	return chatMsgs
}
func DBGroupMessageToChatGroupMessage(dbMsg *db.GroupMessage) *base.GroupMessage {
	if dbMsg == nil {
		return nil
	}
	return &base.GroupMessage{
		FromUserId: dbMsg.FromUserID,
		GroupId:    dbMsg.GroupID,
		Content:    dbMsg.Content,
		CreatedAt:  dbMsg.CreatedAt.Unix(),
	}
}
func DBGroupMessagesToChatGroupMessages(dbMsgs []*db.GroupMessage) []*base.GroupMessage {
	chatMsgs := make([]*base.GroupMessage, 0, len(dbMsgs))
	for _, msg := range dbMsgs {
		chatMsgs = append(chatMsgs, DBGroupMessageToChatGroupMessage(msg))
	}
	return chatMsgs
}
