package db

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
)

type GroupMembers struct {
	GroupID int64
	UserID  int64
	gorm.Model
}

func (GroupMembers) TableName() string {
	return constants.GroupMembersTableName
}

func CheckUserExistInGroup(ctx context.Context, userID, groupID int64) (bool, error) {
	var gm GroupMembers
	err := DB.WithContext(ctx).Where("group_id = ? AND user_id = ?", groupID, userID).First(&gm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CreateGroup(ctx context.Context, userID int64) error {
	// TODO: CreateGroup 逻辑，再创建一个表放群信息¿
	var maxID int64
	err := DB.WithContext(ctx).Model(&GroupMembers{}).Select("MAX(group_id)").Scan(&maxID).Error
	if err != nil {
		return err
	}
	// TODO: 并发Boooooom💥
	groupID := maxID + 1
	err = DB.WithContext(ctx).Create(&GroupMembers{
		GroupID: groupID,
		UserID:  userID,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func AddGroupMember(ctx context.Context, groupID, userID int64) error {
	err := DB.WithContext(ctx).Create(&GroupMembers{
		GroupID: groupID,
		UserID:  userID,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryGroupMemberIDList(ctx context.Context, groupID int64) ([]int64, error) {
	ids := make([]int64, 0)
	err := DB.WithContext(ctx).Model(&GroupMembers{}).Where("group_id = ?", groupID).Pluck("user_id", &ids).Error
	if err != nil {
		return nil, err
	}
	logger.Debugf("Group %d members: %v", groupID, ids)
	return ids, nil
}
