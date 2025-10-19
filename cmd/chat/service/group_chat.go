package service

import (
	"time"

	"github.com/nnieie/golanglab5/cmd/chat/dal/db"
	"github.com/nnieie/golanglab5/cmd/chat/pack"
	"github.com/nnieie/golanglab5/cmd/chat/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *ChatService) SendGroupMessage(groupID, fromUserID int64, content string) error {
	msg := &db.GroupMessage{
		GroupID:    groupID,
		FromUserID: fromUserID,
		Content:    content,
	}
	return db.CreateGroupMessage(msg)
}

func (s *ChatService) GetGroupHistoryMessage(groupID int64, pageNum, pageSize int64) ([]*base.GroupMessage, error) {
	msg, err := db.QueryGroupHistoryMessage(groupID, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return pack.DBGroupMessagesToChatGroupMessages(msg), nil
}

func (s *ChatService) QueryGroupMessageByTime(groupID int64, pageNum, pageSize int64, since int64) ([]*base.GroupMessage, error) {
	msg, err := db.QueryGroupMessageByTime(groupID, pageNum, pageSize, time.Unix(since, 0))
	if err != nil {
		return nil, err
	}
	return pack.DBGroupMessagesToChatGroupMessages(msg), nil
}

func (s *ChatService) GetOfflineGroupMessage(userID, groupID int64, pageNum, pageSize int64) ([]*base.GroupMessage, error) {
	since, err := rpc.QueryUserLastLogoutTime(s.ctx, userID)
	if err != nil {
		return nil, err
	}
	msg, err := db.QueryGroupMessageByTime(groupID, pageNum, pageSize, time.Unix(since, 0))
	if err != nil {
		return nil, err
	}
	return pack.DBGroupMessagesToChatGroupMessages(msg), nil
}
