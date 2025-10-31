package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/nnieie/golanglab5/pkg/logger"
)

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

	cmds, err := pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Errorf("Failed to save group message to Redis: %v", err)
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
