package db

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/pkg/constants"
)

type PrivateMessage struct {
	FromUserID int64
	ToUserID   int64
	Content    string
	gorm.Model
}

func (PrivateMessage) TableName() string {
	return constants.PrivateMessageTableName
}

func CreatePrivateMessage(ctx context.Context, msg *PrivateMessage) error {
	err := DB.WithContext(ctx).Create(msg).Error
	return err
}

func QueryPrivateHistoryMessage(ctx context.Context, fromUser, toUser int64, pageNum, pageSize int64) ([]*PrivateMessage, error) {
	var msgs []*PrivateMessage

	// 使用原生 SQL
	// 用 UNION ALL 解决 from_user_id = ? AND to_user_id = ? OR from_user_id = ? AND to_user_id = ? 使用 OR 索引失效问题
	sql := `
        (SELECT * FROM private_messages 
         WHERE from_user_id = ? AND to_user_id = ? AND deleted_at IS NULL
         ORDER BY created_at DESC)
        UNION ALL
        (SELECT * FROM private_messages 
         WHERE from_user_id = ? AND to_user_id = ? AND deleted_at IS NULL
         ORDER BY created_at DESC)
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `

	offset := (pageNum - 1) * pageSize
	err := DB.WithContext(ctx).Raw(sql,
		fromUser, toUser,
		toUser, fromUser,
		pageSize, offset).Scan(&msgs).Error

	return msgs, err
}

func QueryPrivateMessageByTime(ctx context.Context, fromUser, toUser int64, pageNum, pageSize int64, since time.Time) ([]*PrivateMessage, error) {
	var msgs []*PrivateMessage
	sql := `
        (SELECT * FROM private_messages 
         WHERE from_user_id = ? AND to_user_id = ? AND created_at >= ? AND deleted_at IS NULL
         ORDER BY created_at DESC)
        UNION ALL
        (SELECT * FROM private_messages 
         WHERE from_user_id = ? AND to_user_id = ? AND created_at >= ? AND deleted_at IS NULL
         ORDER BY created_at DESC)
        ORDER BY created_at DESC
        LIMIT ? OFFSET ?
    `

	offset := (pageNum - 1) * pageSize
	err := DB.WithContext(ctx).Raw(sql,
		fromUser, toUser, since,
		toUser, fromUser, since,
		pageSize, offset).Scan(&msgs).Error
	return msgs, err
}

func BatchCreatePrivateMessages(ctx context.Context, messages []*PrivateMessage) error {
	if len(messages) == 0 {
		return nil
	}
	batchSize := 200
	if err := DB.WithContext(ctx).CreateInBatches(messages, batchSize).Error; err != nil {
		return err
	}
	return nil
}
