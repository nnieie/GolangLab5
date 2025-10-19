package db

import (
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

func CheckUserExistInGroup(userID, groupID int64) (bool, error) {
	var gm GroupMembers
	err := DB.Where("group_id = ? AND user_id = ?", groupID, userID).First(&gm).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func CreatGroup(userID int64) error {
	var maxID int64
	err := DB.Model(&GroupMembers{}).Select("MAX(group_id)").Scan(&maxID).Error
	if err != nil {
		return err
	}
	groupID := maxID + 1
	err = DB.Create(&GroupMembers{
		GroupID: groupID,
		UserID:  userID,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func AddGroupMember(groupID, userID int64) error {
	err := DB.Create(&GroupMembers{
		GroupID: groupID,
		UserID:  userID,
	}).Error
	if err != nil {
		return err
	}
	return nil
}

func QueryGroupMemberIDList(groupID int64) ([]int64, error) {
	ids := make([]int64, 0)
	err := DB.Model(&GroupMembers{}).Where("group_id = ?", groupID).Pluck("user_id", &ids).Error
	if err != nil {
		return nil, err
	}
	logger.Debugf("Group %d members: %v", groupID, ids)
	return ids, nil
}
