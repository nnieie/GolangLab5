package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

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
