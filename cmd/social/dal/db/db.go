package db

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/errno"
)

type Follow struct {
	UserID     int64
	FollowerID int64
	gorm.Model
}

func (Follow) TableName() string {
	return constants.FollowsTableName
}

func CreateFollows(ctx context.Context, follow *Follow) error {
	err := DB.WithContext(ctx).Create(follow).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errno.FollowAlreadyExistErr
		}
		return err
	}
	return nil
}

func CheckFollowExist(ctx context.Context, userID, followerID int64) (bool, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&Follow{}).Where("user_id = ? AND follower_id = ?", userID, followerID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func DeleteFollows(ctx context.Context, follow *Follow) error {
	result := DB.WithContext(ctx).Unscoped().Where("user_id = ? AND follower_id = ?", follow.UserID, follow.FollowerID).Delete(&Follow{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func QueryFollowingList(ctx context.Context, userID int64, pageNum, pageSize int64) ([]int64, error) {
	var follows []int64
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}
	err := DB.WithContext(ctx).Model(&Follow{}).Where("follower_id = ?", userID).Order("created_at DESC").
		Limit(int(pageSize)).Offset(int((pageNum-1)*pageSize)).Pluck("user_id", &follows).Error
	if err != nil {
		return nil, err
	}
	return follows, nil
}

func QueryFollowingCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&Follow{}).Where("follower_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func QueryFollowerList(ctx context.Context, userID int64, pageNum, pageSize int64) ([]int64, error) {
	var follows []int64
	err := DB.WithContext(ctx).Model(&Follow{}).Where("user_id = ?", userID).Order("created_at DESC").
		Limit(int(pageSize)).Offset(int((pageNum-1)*pageSize)).Pluck("follower_id", &follows).Error
	if err != nil {
		return nil, err
	}
	return follows, nil
}

func QueryFollowerCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	err := DB.WithContext(ctx).Model(&Follow{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func QueryFriendList(ctx context.Context, userID int64, pageNum, pageSize int64) ([]int64, error) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	// 查询当前用户关注了哪些人
	var followingList []int64
	err := DB.WithContext(ctx).Model(&Follow{}).
		Where("follower_id = ?", userID).
		Pluck("user_id", &followingList).Error
	if err != nil {
		return nil, err
	}
	if len(followingList) == 0 {
		return []int64{}, nil
	}

	// 在这些关注的人中，查询有哪些人也关注了当前用户
	var friends []int64
	err = DB.WithContext(ctx).Model(&Follow{}).
		Where("user_id = ? AND follower_id IN ?", userID, followingList).
		Offset(int((pageNum-1)*pageSize)).
		Limit(int(pageSize)).
		Pluck("follower_id", &friends).Error

	if err != nil {
		return nil, err
	}
	return friends, nil
}

// QueryFriendCount 查询互相关注的好友数量
func QueryFriendCount(ctx context.Context, userID int64) (int64, error) {
	// 查询当前用户关注了哪些人
	var followingList []int64
	err := DB.WithContext(ctx).Model(&Follow{}).
		Where("follower_id = ?", userID).
		Pluck("user_id", &followingList).Error
	if err != nil || len(followingList) == 0 {
		return 0, err
	}

	// 在这些关注的人中，统计有多少人也关注了当前用户
	var count int64
	err = DB.WithContext(ctx).Model(&Follow{}).
		Where("user_id = ? AND follower_id IN ?", userID, followingList).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func CheckFriendExist(ctx context.Context, userID, friendID int64) (bool, error) {
	var count1, count2 int64

	// 检查 A 是否关注了 B
	err := DB.WithContext(ctx).Model(&Follow{}).
		Where("follower_id = ? AND user_id = ?", userID, friendID).
		Count(&count1).Error
	if err != nil || count1 == 0 {
		return false, err
	}

	// 检查 B 是否关注了 A
	err = DB.WithContext(ctx).Model(&Follow{}).
		Where("follower_id = ? AND user_id = ?", friendID, userID).
		Count(&count2).Error
	if err != nil {
		return false, err
	}

	return count2 > 0, nil
}
