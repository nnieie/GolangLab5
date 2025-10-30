package service

import (
	"time"

	"github.com/nnieie/golanglab5/cmd/chat/dal/cache"
	"github.com/nnieie/golanglab5/cmd/chat/dal/db"
	"github.com/nnieie/golanglab5/cmd/chat/pack"
	"github.com/nnieie/golanglab5/cmd/chat/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

// 发送群聊消息
func (s *ChatService) SendGroupMessage(groupID, fromUserID int64, content string) error {
	// 写入Redis缓存
	_, err := cache.SaveGroupMessage(groupID, fromUserID, content)
	if err != nil {
		return err
	}

	return nil
}

// 获取群组历史消息
func (s *ChatService) GetGroupHistoryMessage(groupID int64, pageNum, pageSize int64) ([]*base.GroupMessage, error) {
	// 先尝试从Redis获取
	cachedMsgs, err := cache.GetGroupMessages(groupID, (pageNum-1)*pageSize, pageSize)
	if err == nil && len(cachedMsgs) == int(pageSize) {
		// 转换缓存消息格式
		return pack.ConvertCachedGroupToBaseMessages(cachedMsgs), nil
	}

	// Redis没有数据，从MySQL获取
	msg, err := db.QueryGroupHistoryMessage(s.ctx, groupID, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return pack.DBGroupMessagesToChatGroupMessages(msg), nil
}

func (s *ChatService) QueryGroupMessageByTime(groupID int64, pageNum, pageSize int64, since int64) ([]*base.GroupMessage, error) {
	cacheMsg, err := cache.GetGroupMessagesByTime(groupID, time.Unix(since, 0), (pageNum-1)*pageSize, pageSize)
	if err == nil {
		return pack.ConvertCachedGroupToBaseMessages(cacheMsg), nil
	}

	msg, err := db.QueryGroupMessageByTime(s.ctx, groupID, pageNum, pageSize, time.Unix(since, 0))
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
	msg, err := db.QueryGroupMessageByTime(s.ctx, groupID, pageNum, pageSize, time.Unix(since, 0))
	if err != nil {
		return nil, err
	}
	return pack.DBGroupMessagesToChatGroupMessages(msg), nil
}
