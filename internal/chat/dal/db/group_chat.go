package db

import (
	"context"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/errno"
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

func CreateGroupMessage(ctx context.Context, msg *GroupMessage) error {
	err := DB.WithContext(ctx).Create(msg).Error
	return err
}

func QueryGroupHistoryMessage(ctx context.Context, groupID string, pageNum, pageSize int64) ([]*GroupMessage, error) {
	intGroupID, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil {
		return nil, errno.ParamErr
	}
	var msgs []*GroupMessage
	err = DB.WithContext(ctx).Where("group_id = ?", intGroupID).
		Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Order("created_at DESC, id DESC").Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func QueryGroupMessageByTime(ctx context.Context, groupID string, pageNum, pageSize int64, since time.Time) ([]*GroupMessage, error) {
	intGroupID, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil {
		return nil, errno.ParamErr
	}
	var msgs []*GroupMessage
	err = DB.WithContext(ctx).Where("group_id = ?", intGroupID).
		Where("created_at >= ?", since).
		Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Order("created_at DESC, id DESC").Find(&msgs).Error
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func BatchCreateGroupMessages(ctx context.Context, messages []*GroupMessage) error {
	if len(messages) == 0 {
		return nil
	}
	batchSize := 200
	if err := DB.WithContext(ctx).CreateInBatches(messages, batchSize).Error; err != nil {
		return err
	}

	return nil
}
