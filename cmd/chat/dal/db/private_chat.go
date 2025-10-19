package db

import (
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

func CreatePrivateMessage(msg *PrivateMessage) error {
	err := DB.Create(msg).Error
	return err
}

func QueryPrivateHistoryMessage(fromUser, toUser int64, pageNum, pageSize int64) ([]*PrivateMessage, error) {
	var msgs []*PrivateMessage
	err := DB.Where("from_user_id = ? AND to_user_id = ? OR from_user_id = ? AND to_user_id = ?", fromUser, toUser, toUser, fromUser).
		Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Order("created_at desc").Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func QueryPrivateMessageByTime(fromUser, toUser int64, pageNum, pageSize int64, since time.Time) ([]*PrivateMessage, error) {
	var msgs []*PrivateMessage
	err := DB.Where("from_user_id = ? AND to_user_id = ? OR from_user_id = ? AND to_user_id = ?", fromUser, toUser, toUser, fromUser).
		Where("created_at >= ?", since).
		Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Order("created_at desc").Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
