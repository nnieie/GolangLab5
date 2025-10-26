package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/logger"
)

const (
	privateMessageTTL  = 7 * 24 * time.Hour
	groupMessageTTL    = 7 * 24 * time.Hour
	maxMessageListSize = 1000
)

// 缓存中的私聊消息结构
type CachedPrivateMessage struct {
	ID         string
	FromUserID int64
	ToUserID   int64
	Content    string
	CreatedAt  time.Time
}

// 缓存中的群聊消息结构
type CachedGroupMessage struct {
	ID         string
	GroupID    int64
	FromUserID int64
	Content    string
	CreatedAt  time.Time
}

// Redis键名常量
const (
	PrivateMessagePrefix   = "msg:private:"
	GroupMessagePrefix     = "msg:group:"
	PrivateMessageQueueKey = "msg:queue:private"
	GroupMessageQueueKey   = "msg:queue:group"
)

// 保存私聊消息到Redis
func SavePrivateMessage(ctx context.Context, fromUserID, toUserID int64, content string) (*CachedPrivateMessage, error) {
	// 生成唯一ID
	msgID := fmt.Sprintf("%d_%d_%d", fromUserID, toUserID, time.Now().UnixNano())

	msg := &CachedPrivateMessage{
		ID:         msgID,
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	key := fmt.Sprintf("%s%s", PrivateMessagePrefix, msgID)

	// 管道操作提高性能
	pipe := rChat.Pipeline()

	// 1. 保存消息详情
	pipe.HSet(ctx, key,
		"id", msg.ID,
		"from_user_id", msg.FromUserID,
		"to_user_id", msg.ToUserID,
		"content", msg.Content,
		"created_at", msg.CreatedAt.Format(time.RFC3339Nano),
	)
	pipe.Expire(ctx, key, privateMessageTTL)

	// 2. 添加到用户消息列表（用于快速查询）
	userListKey := fmt.Sprintf("user_msgs:private:%d:%d", max(fromUserID, toUserID), min(fromUserID, toUserID))
	pipe.ZAdd(ctx, userListKey, redis.Z{
		Score:  float64(msg.CreatedAt.UnixNano()),
		Member: msg.ID,
	})
	pipe.ZRemRangeByRank(ctx, userListKey, 0, -maxMessageListSize-1)
	pipe.Expire(ctx, userListKey, privateMessageTTL)

	// 3. 添加到待同步队列
	pipe.LPush(ctx, PrivateMessageQueueKey, msgID)

	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Errorf("Failed to save private message to Redis: %v", err)
		return nil, err
	}

	return msg, nil
}

// 保存群聊消息到Redis
func SaveGroupMessage(groupID, fromUserID int64, content string) (*CachedGroupMessage, error) {
	ctx := context.Background()

	msgID := fmt.Sprintf("%d_%d_%d", groupID, fromUserID, time.Now().UnixNano())

	msg := &CachedGroupMessage{
		ID:         msgID,
		GroupID:    groupID,
		FromUserID: fromUserID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	key := fmt.Sprintf("%s%s", GroupMessagePrefix, msgID)

	pipe := rChat.Pipeline()

	pipe.HSet(ctx, key,
		"id", msg.ID,
		"group_id", msg.GroupID,
		"from_user_id", msg.FromUserID,
		"content", msg.Content,
		"created_at", msg.CreatedAt.Format(time.RFC3339Nano),
	)
	pipe.Expire(ctx, key, groupMessageTTL)

	// 群组消息列表
	groupListKey := fmt.Sprintf("group_msgs:%d", groupID)
	pipe.ZAdd(ctx, groupListKey, redis.Z{
		Score:  float64(msg.CreatedAt.UnixNano()),
		Member: msg.ID,
	})
	pipe.ZRemRangeByRank(ctx, groupListKey, 0, -maxMessageListSize-1)
	pipe.Expire(ctx, groupListKey, groupMessageTTL)

	pipe.LPush(ctx, GroupMessageQueueKey, msgID)

	_, err := pipe.Exec(ctx)
	if err != nil {
		logger.Errorf("Failed to save group message to Redis: %v", err)
		return nil, err
	}

	return msg, nil
}

// 获取用户私聊消息
func GetPrivateMessages(fromUserID, toUserID int64, offset, limit int64) ([]*CachedPrivateMessage, error) {
	ctx := context.Background()

	key := fmt.Sprintf("user_msgs:private:%d:%d", max(fromUserID, toUserID), min(fromUserID, toUserID))

	// 获取消息ID列表（从 Sorted Set 中按分数倒序获取，最新的在前）
	msgIDs, err := rChat.ZRevRange(ctx, key, offset, offset+limit-1).Result()
	if err != nil {
		return nil, err
	}

	if len(msgIDs) == 0 {
		return nil, nil
	}

	// 使用 Pipeline 批量获取消息详情
	pipe := rChat.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(msgIDs))

	for i, msgID := range msgIDs {
		msgKey := fmt.Sprintf("%s%s", PrivateMessagePrefix, msgID)
		cmds[i] = pipe.HGetAll(ctx, msgKey)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var messages []*CachedPrivateMessage
	for _, cmd := range cmds {
		dataMap, err := cmd.Result()
		if err != nil || len(dataMap) == 0 {
			continue
		}

		msg := &CachedPrivateMessage{}
		msg.ID = dataMap["id"]
		msg.Content = dataMap["content"]
		msg.FromUserID, _ = strconv.ParseInt(dataMap["from_user_id"], 10, 64)
		msg.ToUserID, _ = strconv.ParseInt(dataMap["to_user_id"], 10, 64)
		msg.CreatedAt, _ = time.Parse(time.RFC3339Nano, dataMap["created_at"])

		messages = append(messages, msg)
	}

	return messages, nil
}

// 获取群组消息
func GetGroupMessages(groupID int64, offset, limit int64) ([]*CachedGroupMessage, error) {
	ctx := context.Background()

	groupListKey := fmt.Sprintf("group_msgs:%d", groupID)

	// 获取消息ID列表（从 Sorted Set 中按分数倒序获取，最新的在前）
	msgIDs, err := rChat.ZRevRange(ctx, groupListKey, offset, offset+limit-1).Result()
	if err != nil {
		return nil, err
	}

	if len(msgIDs) == 0 {
		return nil, nil
	}

	// 使用 Pipeline 批量获取消息详情
	pipe := rChat.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(msgIDs))

	for i, msgID := range msgIDs {
		msgKey := fmt.Sprintf("%s%s", GroupMessagePrefix, msgID)
		cmds[i] = pipe.HGetAll(ctx, msgKey)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var messages []*CachedGroupMessage
	for _, cmd := range cmds {
		dataMap, err := cmd.Result()
		if err != nil || len(dataMap) == 0 {
			continue
		}

		msg := &CachedGroupMessage{}
		msg.ID = dataMap["id"]
		msg.Content = dataMap["content"]
		msg.FromUserID, _ = strconv.ParseInt(dataMap["from_user_id"], 10, 64)
		msg.GroupID, _ = strconv.ParseInt(dataMap["group_id"], 10, 64)
		msg.CreatedAt, _ = time.Parse(time.RFC3339Nano, dataMap["created_at"])

		messages = append(messages, msg)
	}

	return messages, nil
}

func GetPrivateMessagesByTime(fromUserID, toUserID int64, since time.Time, offset, limit int64) ([]*CachedPrivateMessage, error) {
	ctx := context.Background()

	key := fmt.Sprintf("user_msgs:private:%d:%d", max(fromUserID, toUserID), min(fromUserID, toUserID))

	// 按时间范围查询消息ID
	msgIDs, err := rChat.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%d", since.Unix()),
		Max:    "+inf",
		Offset: offset,
		Count:  limit,
	}).Result()

	if err != nil {
		return nil, err
	}

	if len(msgIDs) == 0 {
		return nil, errno.CacheMissErr
	}

	// 批量获取消息详情
	pipe := rChat.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(msgIDs))

	for i, msgID := range msgIDs {
		msgKey := fmt.Sprintf("%s%s", PrivateMessagePrefix, msgID)
		cmds[i] = pipe.HGetAll(ctx, msgKey)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var messages []*CachedPrivateMessage
	for _, cmd := range cmds {
		dataMap, err := cmd.Result()
		if err != nil || len(dataMap) == 0 {
			continue
		}

		msg := &CachedPrivateMessage{}
		msg.ID = dataMap["id"]
		msg.Content = dataMap["content"]
		msg.FromUserID, _ = strconv.ParseInt(dataMap["from_user_id"], 10, 64)
		msg.ToUserID, _ = strconv.ParseInt(dataMap["to_user_id"], 10, 64)
		msg.CreatedAt, _ = time.Parse(time.RFC3339Nano, dataMap["created_at"])

		messages = append(messages, msg)
	}

	return messages, nil
}

func GetGroupMessagesByTime(groupID int64, since time.Time, offset, limit int64) ([]*CachedGroupMessage, error) {
	ctx := context.Background()

	key := fmt.Sprintf("user_msgs:group:%d", groupID)

	// 按时间范围查询消息ID
	msgIDs, err := rChat.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    fmt.Sprintf("%d", since.Unix()),
		Max:    "+inf",
		Offset: offset,
		Count:  limit,
	}).Result()

	if err != nil {
		return nil, err
	}

	if len(msgIDs) == 0 {
		return nil, errno.CacheMissErr
	}

	// 批量获取消息详情
	pipe := rChat.Pipeline()
	cmds := make([]*redis.MapStringStringCmd, len(msgIDs))

	for i, msgID := range msgIDs {
		msgKey := fmt.Sprintf("%s%s", GroupMessagePrefix, msgID)
		cmds[i] = pipe.HGetAll(ctx, msgKey)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return nil, err
	}

	// 解析结果
	var messages []*CachedGroupMessage
	for _, cmd := range cmds {
		dataMap, err := cmd.Result()
		if err != nil || len(dataMap) == 0 {
			continue
		}

		msg := &CachedGroupMessage{}
		msg.ID = dataMap["id"]
		msg.Content = dataMap["content"]
		msg.FromUserID, _ = strconv.ParseInt(dataMap["from_user_id"], 10, 64)
		msg.GroupID, _ = strconv.ParseInt(dataMap["group_id"], 10, 64)
		msg.CreatedAt, _ = time.Parse(time.RFC3339Nano, dataMap["created_at"])

		messages = append(messages, msg)
	}

	return messages, nil
}
