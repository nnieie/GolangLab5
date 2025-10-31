package db

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/cmd/interaction/rpc"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/errno"
)

const (
	VideoLikeType   = 1
	CommentLikeType = 2

	LikeActionType   = 1
	UnlikeActionType = 2
)

type Like struct {
	UserID   int64
	TargetID int64
	Type     int64
	gorm.Model
}

func (Like) TableName() string {
	return constants.LikeTableName
}

func LikeAction(ctx context.Context, like *Like) (int64, error) {
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(like).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				return errno.LikeAlreadyExistErr
			}
			return err
		}
		// 更新点赞数
		var updateErr error
		switch like.Type {
		case VideoLikeType:
			_, updateErr = rpc.UpdateVideoLikeCount(ctx, like.TargetID, 1)
		case CommentLikeType:
			updateErr = tx.Model(&Comment{}).Where("id = ?", like.TargetID).Update("like_count", gorm.Expr("like_count + 1")).Error
		}

		if updateErr != nil {
			return updateErr
		}

		return nil
	})

	if err != nil {
		return 0, err
	}
	return int64(like.ID), nil
}

func UnlikeAction(ctx context.Context, userID, targetID int64, likeType int64) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var result *gorm.DB
		switch likeType {
		case VideoLikeType:
			result = tx.Where("user_id = ? AND target_id = ? AND type = ?", userID, targetID, VideoLikeType).Delete(&Like{})
		case CommentLikeType:
			result = tx.Where("user_id = ? AND target_id = ? AND type = ?", userID, targetID, CommentLikeType).Delete(&Like{})
		default:
			return errno.ParamErr
		}

		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errno.LikeIsNotExistErr
		}

		var err error
		switch likeType {
		case VideoLikeType:
			_, err = rpc.UpdateVideoLikeCount(ctx, targetID, -1)
		case CommentLikeType:
			err = tx.Model(&Comment{}).Where("id = ?", targetID).Update("like_count", gorm.Expr("like_count - 1")).Error
		}
		if err != nil {
			return err
		}
		return nil
	})
}

func UnlikeActionByID(ctx context.Context, id int64, likeType int64) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var likeToDelete Like
		// 查询要删除的 Like 记录
		if err := tx.Where("id = ?", id).First(&likeToDelete).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errno.LikeIsNotExistErr
			}
			return err
		}

		// 根据 ID 删除
		if err := tx.Delete(&Like{}, id).Error; err != nil {
			return err
		}

		// 更新计数器
		var err error
		switch likeType {
		case VideoLikeType:
			_, err = rpc.UpdateVideoLikeCount(ctx, likeToDelete.TargetID, -1)
		case CommentLikeType:
			err = tx.Model(&Comment{}).Where("id = ?", likeToDelete.TargetID).Update("like_count", gorm.Expr("like_count - 1")).Error
		}

		return err
	})
}

func QueryLikeVideoListByUserID(ctx context.Context, userID int64, pageNum, pageSize int64) ([]int64, error) {
	var likes []int64
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	err := DB.WithContext(ctx).Model(&Like{}).Where("type = 1 AND user_id = ?", userID).Order("created_at DESC").Pluck("target_id", &likes).
		Offset(int((pageNum - 1) * pageSize)).Limit(int(pageSize)).Error
	if err != nil {
		return nil, err
	}
	return likes, nil
}

func QueryLikeByUserIDAndTargetIDAndType(ctx context.Context, userID, targetID int64, likeType int64) (*Like, error) {
	var like Like
	err := DB.WithContext(ctx).Where("user_id = ? AND target_id = ? AND type = ?", userID, targetID, likeType).Order("created_at DESC").First(&like).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &like, nil
}
