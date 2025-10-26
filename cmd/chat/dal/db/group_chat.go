package db

import (
	"time"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/pkg/constants"
)

type GroupMessage struct {
	FromUserID int64
	GroupID    int64
	Content    string
	gorm.Model
}

func (GroupMessage) TableName() string {
	return constants.GroupMessageTableName
}

func CreateGroupMessage(msg *GroupMessage) error {
	err := DB.Create(msg).Error
	return err
}

func QueryGroupHistoryMessage(groupID int64, pageNum, pageSize int64) ([]*GroupMessage, error) {
	var msgs []*GroupMessage
	err := DB.Where("group_id = ?", groupID).
		Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Order("created_at desc").Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func QueryGroupMessageByTime(groupID int64, pageNum, pageSize int64, since time.Time) ([]*GroupMessage, error) {
	var msgs []*GroupMessage
	err := DB.Where("group_id = ?", groupID).
		Where("created_at >= ?", since).
		Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Order("created_at desc").Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func BatchCreateGroupMessages(messages []*GroupMessage) error {
	if len(messages) == 0 {
		return nil
	}
	if err := DB.CreateInBatches(messages, len(messages)).Error; err != nil {
		return err
	}

	return nil
}
