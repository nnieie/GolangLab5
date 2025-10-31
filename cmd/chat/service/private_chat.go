package service

import (
	"time"

	"github.com/nnieie/golanglab5/cmd/chat/dal/cache"
	"github.com/nnieie/golanglab5/cmd/chat/dal/db"
	"github.com/nnieie/golanglab5/cmd/chat/pack"
	"github.com/nnieie/golanglab5/cmd/chat/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

// 发送私聊消息
func (s *ChatService) SendPrivateMessage(fromUser, toUser int64, content string) error {
	// 写入Redis缓存
	_, err := cache.SavePrivateMessage(s.ctx, fromUser, toUser, content)
	if err != nil {
		return err
	}

	return nil
}

// 获取历史消息
func (s *ChatService) GetPrivateHistoryMessage(fromUser, toUser int64, pageNum, pageSize int64) ([]*base.PrivateMessage, error) {
	// 先尝试从Redis获取
	cachedMsgs, err := cache.GetPrivateMessages(fromUser, toUser, (pageNum-1)*pageSize, pageSize)
	if err == nil && len(cachedMsgs) == int(pageSize) {
		// 转换缓存消息格式
		return pack.ConvertCachedToBaseMessages(cachedMsgs), nil
	}

	// Redis没有数据，从MySQL获取
	msg, err := db.QueryPrivateHistoryMessage(s.ctx, fromUser, toUser, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return pack.DBPrivateMessagesToChatPrivateMessages(msg), nil
}

func (s *ChatService) GetPrivateMessageByTime(fromUser, toUser int64, pageNum, pageSize int64, since int64) ([]*base.PrivateMessage, error) {
	cacheMsg, err := cache.GetPrivateMessagesByTime(fromUser, toUser, time.Unix(since, 0), (pageNum-1)*pageSize, pageSize)
	if err == nil {
		return pack.ConvertCachedToBaseMessages(cacheMsg), nil
	}

	msg, err := db.QueryPrivateMessageByTime(s.ctx, fromUser, toUser, pageNum, pageSize, time.Unix(since, 0))
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
	msg, err := db.QueryPrivateMessageByTime(s.ctx, fromUser, toUser, pageNum, pageSize, time.Unix(since, 0))
	if err != nil {
		return nil, err
	}
	return pack.DBPrivateMessagesToChatPrivateMessages(msg), nil
}
