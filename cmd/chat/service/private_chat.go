package service

import (
	"strconv"
	"time"

	"github.com/nnieie/golanglab5/cmd/chat/dal/cache"
	"github.com/nnieie/golanglab5/cmd/chat/dal/db"
	"github.com/nnieie/golanglab5/cmd/chat/pack"
	"github.com/nnieie/golanglab5/cmd/chat/rpc"
	"github.com/nnieie/golanglab5/kitex_gen/base"
	"github.com/nnieie/golanglab5/pkg/tracer"
)

// 发送私聊消息
func (s *ChatService) SendPrivateMessage(fromUser, toUser string, content string) error {
	fromUserID, err := strconv.ParseInt(fromUser, 10, 64)
	if err != nil {
		return err
	}
	toUserID, err := strconv.ParseInt(toUser, 10, 64)
	if err != nil {
		return err
	}
	// 写入Redis缓存
	_, err = cache.SavePrivateMessage(s.ctx, fromUserID, toUserID, content)
	if err != nil {
		return err
	}

	tracer.ChatMessageCounter.Add(s.ctx, 1)
	return nil
}

// 获取历史消息
func (s *ChatService) GetPrivateHistoryMessage(fromUser, toUser string, pageNum, pageSize int64) ([]*base.PrivateMessage, error) {
	fromUserID, err := strconv.ParseInt(fromUser, 10, 64)
	if err != nil {
		return nil, err
	}
	toUserID, err := strconv.ParseInt(toUser, 10, 64)
	if err != nil {
		return nil, err
	}
	// 先尝试从Redis获取
	cachedMsgs, err := cache.GetPrivateMessages(fromUserID, toUserID, (pageNum-1)*pageSize, pageSize)
	if err == nil && len(cachedMsgs) == int(pageSize) {
		// 转换缓存消息格式
		return pack.ConvertCachedToBaseMessages(cachedMsgs), nil
	}

	// Redis没有数据，从MySQL获取
	msg, err := db.QueryPrivateHistoryMessage(s.ctx, fromUserID, toUserID, pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	return pack.DBPrivateMessagesToChatPrivateMessages(msg), nil
}

func (s *ChatService) GetPrivateMessageByTime(fromUser, toUser string, pageNum, pageSize int64, since int64) ([]*base.PrivateMessage, error) {
	fromUserID, err := strconv.ParseInt(fromUser, 10, 64)
	if err != nil {
		return nil, err
	}
	toUserID, err := strconv.ParseInt(toUser, 10, 64)
	if err != nil {
		return nil, err
	}
	cacheMsg, err := cache.GetPrivateMessagesByTime(fromUserID, toUserID, time.Unix(since, 0), (pageNum-1)*pageSize, pageSize)
	if err == nil {
		return pack.ConvertCachedToBaseMessages(cacheMsg), nil
	}

	msg, err := db.QueryPrivateMessageByTime(s.ctx, fromUserID, toUserID, pageNum, pageSize, time.Unix(since, 0))
	if err != nil {
		return nil, err
	}
	return pack.DBPrivateMessagesToChatPrivateMessages(msg), nil
}

func (s *ChatService) GetOfflinePrivateMessage(fromUser, toUser string, pageNum, pageSize int64) ([]*base.PrivateMessage, error) {
	fromUserID, err := strconv.ParseInt(fromUser, 10, 64)
	if err != nil {
		return nil, err
	}
	toUserID, err := strconv.ParseInt(toUser, 10, 64)
	if err != nil {
		return nil, err
	}
	since, err := rpc.QueryUserLastLogoutTime(s.ctx, fromUser)
	if err != nil {
		return nil, err
	}
	msg, err := db.QueryPrivateMessageByTime(s.ctx, fromUserID, toUserID, pageNum, pageSize, time.Unix(since, 0))
	if err != nil {
		return nil, err
	}
	return pack.DBPrivateMessagesToChatPrivateMessages(msg), nil
}
