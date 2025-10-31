package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
		if err := w.syncPrivateMessages(); err != nil {
			logger.Errorf("sync private messages failed: %v", err)
		}
		if err := w.syncGroupMessages(); err != nil {
			logger.Errorf("sync group messages failed: %v", err)
		}
	}
}

// 同步私聊消息到MySQL
func (w *SyncWorker) syncPrivateMessages() error {
	ctx := context.Background()

	// 批量获取待同步消息
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

	// 处理每条消息
	for _, msgID := range privateMsgIDs {
		if msg, err := w.getPrivateMessageFromCache(msgID); err == nil {
			privateMessages = append(privateMessages, &db.PrivateMessage{
				FromUserID: msg.FromUserID,
				ToUserID:   msg.ToUserID,
				Content:    msg.Content,
				Model: gorm.Model{
					CreatedAt: msg.CreatedAt,
				},
			})
			processedMsgs = append(processedMsgs, msgID)
		} else {
			logger.Warnf("Failed to get private message %s from cache: %v", msgID, err)
		}
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

	// 记录同步指标
	logger.Infof("successfully synced %d private messages to MySQL", syncCount)
	return nil
}

// 同步群聊消息到MySQL
func (w *SyncWorker) syncGroupMessages() error {
	ctx := context.Background()

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

	for _, msgID := range groupMsgIDs {
		if msg, err := w.getGroupMessageFromCache(msgID); err == nil {
			groupMessages = append(groupMessages, &db.GroupMessage{
				FromUserID: msg.FromUserID,
				GroupID:    msg.GroupID,
				Content:    msg.Content,
				Model: gorm.Model{
					CreatedAt: msg.CreatedAt,
				},
			})
			processedMsgs = append(processedMsgs, msgID)
		} else {
			logger.Warnf("Failed to get group message %s from cache: %v", msgID, err)
		}
	}

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
			pipe.LRem(ctx, PrivateMessageQueueKey, 1, msgID)
		}
		if _, err := pipe.Exec(ctx); err != nil {
			logger.Errorf("Failed to remove processed messages from queue: %v", err)
		}
	}

	// 记录同步指标
	logger.Infof("successfully synced %d group messages to MySQL", syncCount)

	return nil
}

// 从缓存获取私聊消息
func (w *SyncWorker) getPrivateMessageFromCache(msgID string) (*CachedPrivateMessage, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", PrivateMessagePrefix, msgID)

	msgData, err := rChat.HGet(ctx, key, "data").Result()
	if err != nil {
		return nil, err
	}

	var msg CachedPrivateMessage
	if err := json.Unmarshal([]byte(msgData), &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

// 从缓存获取群聊消息
func (w *SyncWorker) getGroupMessageFromCache(msgID string) (*CachedGroupMessage, error) {
	ctx := context.Background()
	key := fmt.Sprintf("%s%s", GroupMessagePrefix, msgID)

	msgData, err := rChat.HGet(ctx, key, "data").Result()
	if err != nil {
		return nil, err
	}

	var msg CachedGroupMessage
	if err := json.Unmarshal([]byte(msgData), &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}
