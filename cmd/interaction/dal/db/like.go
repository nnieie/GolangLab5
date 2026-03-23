package db

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/nnieie/golanglab5/cmd/interaction/rpc"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/logger"
)

const (
	VideoLikeType   = 1
	CommentLikeType = 2

	LikeActionType   = 1
	UnlikeActionType = 2

	batchSize = 100
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
		// 如果不存在 -> 插入
		// 如果存在但 deleted_at IS NULL -> 不做任何事
		// 如果存在且 deleted_at IS NOT NULL -> 更新 deleted_at = NULL
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}, {Name: "target_id"}, {Name: "type"}},
			DoUpdates: clause.Assignments(map[string]interface{}{"deleted_at": nil}),
		}).Create(like).Error; err != nil {
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

	offset := int((pageNum - 1) * pageSize)
	limit := int(pageSize)

	err := DB.WithContext(ctx).Model(&Like{}).
		Where("type = ? AND user_id = ?", VideoLikeType, userID).
		Order("created_at DESC, id DESC").
		Offset(offset).
		Limit(limit).
		Pluck("target_id", &likes).Error
	if err != nil {
		return nil, err
	}
	return likes, nil
}

func QueryLikeByUserIDAndTargetIDAndType(ctx context.Context, userID, targetID int64, likeType int64) (*Like, error) {
	var like Like
	err := DB.WithContext(ctx).Where("user_id = ? AND target_id = ? AND type = ?", userID, targetID, likeType).First(&like).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &like, nil
}

// BatchLikeAction 批量点赞操作
func BatchLikeAction(ctx context.Context, likes []Like) error {
	if len(likes) == 0 {
		return nil
	}
	videoLikeCounts := make(map[int64]int64)
	commentLikeCounts := make(map[int64]int64)

	for _, like := range likes {
		switch like.Type {
		case VideoLikeType:
			videoLikeCounts[like.TargetID]++
		case CommentLikeType:
			commentLikeCounts[like.TargetID]++
		}
	}

	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 批量 upsert
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "user_id"}, {Name: "target_id"}, {Name: "type"}},
			// 冲突时把 deleted_at 设回 NULL
			DoUpdates: clause.Assignments(map[string]interface{}{"deleted_at": nil}),
		}).CreateInBatches(likes, batchSize).Error; err != nil {
			return err
		}

		// 批量更新评论点赞数
		if len(commentLikeCounts) > 0 {
			query := "UPDATE comments SET like_count = like_count + CASE id "
			var args []interface{}
			var ids []int64

			for id, count := range commentLikeCounts {
				query += "WHEN ? THEN ? "
				args = append(args, id, count)
				ids = append(ids, id)
			}
			query += "END WHERE id IN ?"
			args = append(args, ids)

			if err := tx.Exec(query, args...).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 调用 Video 服务批量增加点赞数
	if len(videoLikeCounts) > 0 {
		go func() {
			if _, err := rpc.BatchUpdateVideoLikeCount(context.Background(), videoLikeCounts); err != nil {
				logger.Errorf("Failed to async batch update video like count: %v", err)
			}
		}()
	}

	return nil
}

// BatchUnlikeAction 批量取消点赞操作
func BatchUnlikeAction(ctx context.Context, unlikes []Like) error {
	if len(unlikes) == 0 {
		return nil
	}

	type targetKey struct {
		TargetID int64
		Type     int64
	}

	deleteGroups := make(map[targetKey][]int64)
	videoLikeCounts := make(map[int64]int64)
	commentLikeCounts := make(map[int64]int64)

	for _, unlike := range unlikes {
		key := targetKey{TargetID: unlike.TargetID, Type: unlike.Type}
		deleteGroups[key] = append(deleteGroups[key], unlike.UserID)
	}

	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 批量执行 Delete
		for key, userIDs := range deleteGroups {
			result := tx.Where("target_id = ? AND type = ? AND user_id IN ?",
				key.TargetID, key.Type, userIDs).Delete(&Like{})
			if result.Error != nil {
				return result.Error
			}
			if result.RowsAffected == 0 {
				continue
			}

			switch key.Type {
			case VideoLikeType:
				videoLikeCounts[key.TargetID] += result.RowsAffected
			case CommentLikeType:
				commentLikeCounts[key.TargetID] += result.RowsAffected
			}
		}

		// 批量更新评论点赞数
		if len(commentLikeCounts) > 0 {
			query := "UPDATE comments SET like_count = like_count - CASE id "
			var args []interface{}
			var ids []int64

			for id, count := range commentLikeCounts {
				query += "WHEN ? THEN ? "
				args = append(args, id, count)
				ids = append(ids, id)
			}
			query += "END WHERE id IN ?"
			args = append(args, ids)

			if err := tx.Exec(query, args...).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	// 调用 Video 服务批量扣减点赞数
	if len(videoLikeCounts) > 0 {
		go func() {
			decrementCounts := make(map[int64]int64, len(videoLikeCounts))
			for vid, count := range videoLikeCounts {
				decrementCounts[vid] = -count
			}
			if _, err := rpc.BatchUpdateVideoLikeCount(context.Background(), decrementCounts); err != nil {
				logger.Errorf("Failed to async batch update video like count (unlike): %v", err)
			}
		}()
	}

	return nil
}
