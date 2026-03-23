package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/nnieie/golanglab5/pkg/errno"
)

// 按时间范围获取群组消息
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
