package cache

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/cmd/chat/dal/db"
	"github.com/nnieie/golanglab5/pkg/logger"
)

const (
	batchSize    = 200             // 每批处理200条消息
	syncInterval = 5 * time.Second // 每5秒同步一次
)

type SyncWorker struct {
	batchSize    int
	syncInterval time.Duration
}

func NewSyncWorker() *SyncWorker {
	return &SyncWorker{
		batchSize:    batchSize,
		syncInterval: syncInterval,
	}
}

// 启动同步工作器
func (w *SyncWorker) Start() {
	ticker := time.NewTicker(w.syncInterval)
	defer ticker.Stop()

	for range ticker.C {
		eg := &errgroup.Group{}
		eg.Go(func() error {
			if err := w.syncPrivateMessages(); err != nil {
				logger.Errorf("sync private messages failed: %v", err)
				return err
			}
			return nil
		})

		eg.Go(func() error {
			if err := w.syncGroupMessages(); err != nil {
				logger.Errorf("sync group messages failed: %v", err)
				return err
			}
			return nil
		})
		if err := eg.Wait(); err != nil {
			logger.Errorf("sync messages failed: %v", err)
		}
	}
}

// 同步私聊消息到MySQL
func (w *SyncWorker) syncPrivateMessages() error {
	ctx := context.Background()

	// 批量获取待同步消息ID
	privateMsgIDs, err := rChat.LRange(ctx, PrivateMessageQueueKey, 0, int64(w.batchSize-1)).Result()
	if err != nil {
		logger.Errorf("sync message to mysql err: %v", err)
		return err
	}

	if len(privateMsgIDs) == 0 {
		return nil
	}

	logger.Infof("Starting to sync %d private messages...", len(privateMsgIDs))

	var privateMessages []*db.PrivateMessage
	var processedMsgs []string

	// 使用 Pipeline 批量获取所有消息
	pipe := rChat.Pipeline()
	cmds := make(map[string]*redis.MapStringStringCmd, len(privateMsgIDs))

	for _, msgID := range privateMsgIDs {
		key := fmt.Sprintf("%s%s", PrivateMessagePrefix, msgID)
		cmds[msgID] = pipe.HGetAll(ctx, key) // 使用 HGetAll 获取所有字段
	}

	// 一次性执行所有 Redis 命令
	_, err = pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Errorf("Failed to execute pipeline: %v", err)
		return err
	}

	// 遍历所有命令结果
	for msgID, cmd := range cmds {
		msgData, err := cmd.Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				logger.Warnf("Message %s not found in cache", msgID)
			} else {
				logger.Warnf("Failed to get message %s from cache: %v", msgID, err)
			}
			continue
		}

		// msgData 是 map[string]string
		if len(msgData) == 0 {
			logger.Warnf("Message %s is empty", msgID)
			continue
		}

		// 解析字段
		fromUserID, _ := strconv.ParseInt(msgData["from_user_id"], 10, 64)
		toUserID, _ := strconv.ParseInt(msgData["to_user_id"], 10, 64)
		content := msgData["content"]
		createdAtStr := msgData["created_at"]

		// 解析时间
		createdAt, err := time.Parse(time.RFC3339Nano, createdAtStr)
		if err != nil {
			logger.Warnf("Failed to parse created_at for message %s: %v", msgID, err)
			createdAt = time.Now()
		}

		// 组装数据
		privateMessages = append(privateMessages, &db.PrivateMessage{
			FromUserID: fromUserID,
			ToUserID:   toUserID,
			Content:    content,
			Model: gorm.Model{
				CreatedAt: createdAt,
			},
		})
		processedMsgs = append(processedMsgs, msgID)
	}

	// 批量写入MySQL
	syncCount := 0
	if len(privateMessages) > 0 {
		if err := db.BatchCreatePrivateMessages(ctx, privateMessages); err != nil {
			return err
		}
		syncCount += len(privateMessages)
		logger.Infof("Synced %d private messages to MySQL", len(privateMessages))
	}

	// 从队列中移除已处理的消息
	if len(processedMsgs) > 0 {
		pipe := rChat.Pipeline()
		for _, msgID := range processedMsgs {
			pipe.LRem(ctx, PrivateMessageQueueKey, 1, msgID)
		}
		if _, err := pipe.Exec(ctx); err != nil {
			logger.Errorf("Failed to remove processed messages from queue: %v", err)
		}
	}

	logger.Infof("successfully synced %d private messages to MySQL", syncCount)
	return nil
}

// 同步群聊消息到MySQL
func (w *SyncWorker) syncGroupMessages() error {
	ctx := context.Background()

	// 批量获取待同步消息ID
	groupMsgIDs, err := rChat.LRange(ctx, GroupMessageQueueKey, 0, int64(w.batchSize-1)).Result()
	if err != nil {
		logger.Errorf("sync message to mysql err: %v", err)
		return err
	}

	if len(groupMsgIDs) == 0 {
		return nil
	}

	logger.Infof("Starting to sync %d group messages...", len(groupMsgIDs))

	var groupMessages []*db.GroupMessage
	var processedMsgs []string

	// 使用 Pipeline 批量获取所有消息
	pipe := rChat.Pipeline()
	cmds := make(map[string]*redis.MapStringStringCmd, len(groupMsgIDs))

	for _, msgID := range groupMsgIDs {
		key := fmt.Sprintf("%s%s", GroupMessagePrefix, msgID)
		cmds[msgID] = pipe.HGetAll(ctx, key) // 使用 HGetAll 获取所有字段
	}

	// 一次性执行所有 Redis 命令
	_, err = pipe.Exec(ctx)
	if err != nil && !errors.Is(err, redis.Nil) {
		logger.Errorf("Failed to execute pipeline: %v", err)
		return err
	}

	// 遍历所有命令结果
	for msgID, cmd := range cmds {
		msgData, err := cmd.Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				logger.Warnf("Message %s not found in cache", msgID)
			} else {
				logger.Warnf("Failed to get message %s from cache: %v", msgID, err)
			}
			continue
		}

		// msgData 是 map[string]string
		if len(msgData) == 0 {
			logger.Warnf("Message %s is empty", msgID)
			continue
		}

		// 解析字段
		fromUserID, _ := strconv.ParseInt(msgData["from_user_id"], 10, 64)
		groupID, _ := strconv.ParseInt(msgData["group_id"], 10, 64)
		content := msgData["content"]
		createdAtStr := msgData["created_at"]

		// 解析时间
		createdAt, err := time.Parse(time.RFC3339Nano, createdAtStr)
		if err != nil {
			logger.Warnf("Failed to parse created_at for message %s: %v", msgID, err)
			createdAt = time.Now()
		}

		// 组装数据
		groupMessages = append(groupMessages, &db.GroupMessage{
			FromUserID: fromUserID,
			GroupID:    groupID,
			Content:    content,
			Model: gorm.Model{
				CreatedAt: createdAt,
			},
		})
		processedMsgs = append(processedMsgs, msgID)
	}

	// 批量写入MySQL
	syncCount := 0
	if len(groupMessages) > 0 {
		if err := db.BatchCreateGroupMessages(ctx, groupMessages); err != nil {
			return err
		}
		syncCount += len(groupMessages)
		logger.Infof("Synced %d group messages to MySQL", len(groupMessages))
	}

	// 从队列中移除已处理的消息
	if len(processedMsgs) > 0 {
		pipe := rChat.Pipeline()
		for _, msgID := range processedMsgs {
			pipe.LRem(ctx, GroupMessageQueueKey, 1, msgID)
		}
		if _, err := pipe.Exec(ctx); err != nil {
			logger.Errorf("Failed to remove processed messages from queue: %v", err)
		}
	}

	logger.Infof("successfully synced %d group messages to MySQL", syncCount)
	return nil
}
