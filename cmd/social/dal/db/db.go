package db

import (
	"context"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/errno"
)

type Follow struct {
	UserID     int64
	FollowerID int64
	gorm.Model
}

func CreateFollows(ctx context.Context, follow *Follow) error {
	FollowExist, err := CheckFollowExist(ctx, follow.UserID, follow.FollowerID)
	if err != nil {
		return err
	}
	if FollowExist {
		return errno.FollowAlreadyExistErr
	}
	err = DB.WithContext(ctx).Create(follow).Error
	if err != nil {
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
	result := DB.WithContext(ctx).Model(&Follow{}).Where("user_id = ? AND follower_id = ?", follow.UserID, follow.FollowerID).Delete(&Follow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errno.FollowIsNotExistErr
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
	err := DB.WithContext(ctx).Model(&Follow{}).Where("follower_id = ?", userID).Limit(int(pageSize)).Offset(int(pageNum)-1).Pluck("user_id", &follows).Error
	if err != nil {
		return nil, err
	}
	if len(follows) == 0 {
		return nil, errno.FollowIsNotExistErr
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
	err := DB.WithContext(ctx).Model(&Follow{}).Where("user_id = ?", userID).Limit(int(pageSize)).Offset(int(pageNum)-1).Pluck("follower_id", &follows).Error
	if err != nil {
		return nil, err
	}
	if len(follows) == 0 {
		return nil, errno.FollowIsNotExistErr
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
	var friends []int64
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 20
	}

	// 使用 JOIN 查询互相关注的用户
	// t1 代表 user_id 关注的人的记录 (A -> B)
	// t2 代表关注 user_id 的人的记录 (B -> A)
	// 通过 JOIN 将这两个关系连接起来，找到互相关注的记录
	err := DB.WithContext(ctx).Table(constants.FollowsTableName+" as t1").
		Select("t1.user_id").
		Joins("JOIN "+constants.FollowsTableName+" as t2 ON t1.user_id = t2.follower_id AND t1.follower_id = t2.user_id").
		Where("t1.follower_id = ? AND t1.deleted_at IS NULL AND t2.deleted_at IS NULL", userID).
		Offset(int(pageNum)-1).
		Limit(int(pageSize)).
		Pluck("user_id", &friends).Error

	if err != nil {
		return nil, err
	}
	if len(friends) == 0 {
		return nil, errno.FriendIsNotExistErr
	}
	return friends, nil
}

// QueryFriendCount 查询互相关注的好友数量
func QueryFriendCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	// 使用与 QueryFriendList 相同的 JOIN 逻辑来计数
	err := DB.WithContext(ctx).
		Table(constants.FollowsTableName+" as t1").
		Joins("JOIN "+constants.FollowsTableName+" as t2 ON t1.user_id = t2.follower_id AND t1.follower_id = t2.user_id").
		Where("t1.follower_id = ? AND t1.deleted_at IS NULL AND t2.deleted_at IS NULL", userID).
		Count(&count).Error

	if err != nil {
		return 0, err
	}
	return count, nil
}

func CheckFriendExist(ctx context.Context, userID, friendID int64) (bool, error) {
	var count int64
	// 检查互相关注：user_id 关注 friend_id，且 friend_id 关注 user_id
	err := DB.WithContext(ctx).
		Table(constants.FollowsTableName+" as t1").
		Joins("JOIN "+constants.FollowsTableName+" as t2 ON t1.user_id = t2.follower_id AND t1.follower_id = t2.user_id").
		Where("t1.follower_id = ? AND t1.user_id = ? AND t1.deleted_at IS NULL AND t2.deleted_at IS NULL", userID, friendID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil // 如果 count > 0，表示互相关注
}
