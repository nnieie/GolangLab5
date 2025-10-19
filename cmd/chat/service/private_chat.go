package service

import (
	"time"

	"github.com/nnieie/golanglab5/cmd/chat/dal/db"
	"github.com/nnieie/golanglab5/cmd/chat/pack"
	"github.com/nnieie/golanglab5/cmd/chat/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *ChatService) SendPrivateMessage(fromUser, toUser int64, content string) error {
	msg := &db.PrivateMessage{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Content:    content,
	}
	return db.CreatePrivateMessage(msg)
}

func (s *ChatService) GetPrivateHistoryMessage(fromUser, toUser int64, pageNum, pageSize int64) ([]*base.PrivateMessage, error) {
	msg, err := db.QueryPrivateHistoryMessage(fromUser, toUser, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return pack.DBPrivateMessagesToChatPrivateMessages(msg), nil
}

func (s *ChatService) GetPrivateMessageByTime(fromUser, toUser int64, pageNum, pageSize int64, since int64) ([]*base.PrivateMessage, error) {
	msg, err := db.QueryPrivateMessageByTime(fromUser, toUser, pageNum, pageSize, time.Unix(since, 0))
	if err != nil {
		return nil, err
	}
	return pack.DBPrivateMessagesToChatPrivateMessages(msg), nil
}

func (s *ChatService) GetOfflinePrivateMessage(fromUser, toUser int64, pageNum, pageSize int64) ([]*base.PrivateMessage, error) {
	since, err := rpc.QueryUserLastLogoutTime(s.ctx, fromUser)
	if err != nil {
		return nil, err
	}
	msg, err := db.QueryPrivateMessageByTime(fromUser, toUser, pageNum, pageSize, time.Unix(since, 0))
	if err != nil {
		return nil, err
	}
	return pack.DBPrivateMessagesToChatPrivateMessages(msg), nil
}
