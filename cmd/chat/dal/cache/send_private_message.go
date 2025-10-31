package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/nnieie/golanglab5/pkg/logger"
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

	cmds, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Errorf("Failed to save private message to Redis: %v", err)
		return nil, err
	}

	// 检查每个命令的执行结果
	for i, cmd := range cmds {
		if cmd.Err() != nil && !errors.Is(cmd.Err(), redis.Nil) {
			logger.Warnf("Pipeline command %d failed: %v", i, cmd.Err())
		}
	}

	return msg, nil
}
