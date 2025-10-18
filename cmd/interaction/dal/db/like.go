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
	switch like.Type {
	case VideoLikeType:
		if video, err := rpc.QueryVideoByID(ctx, like.TargetID); err != nil {
			return 0, err
		} else if video == nil {
			return 0, errno.VideoIsNotExistErr
		}
	case CommentLikeType:
		if _, err := QueryCommentByID(ctx, like.TargetID); err != nil {
			return 0, err
		}
	}
	dbLike, err := QueryLikeByUserIDAndTargetIDAndType(ctx, like.UserID, like.TargetID, like.Type)
	if err != nil {
		return 0, err
	}
	if dbLike != nil {
		return 0, errno.LikeAlreadyExistErr
	}

	err = DB.Create(like).Error
	if err != nil {
		return 0, err
	}

	switch like.Type {
	case VideoLikeType:
		_, err = rpc.UpdateVideoLikeCount(ctx, like.TargetID, 1)
	case CommentLikeType:
		err = DB.Model(&Comment{}).Where("id = ?", like.TargetID).Update("like_count", gorm.Expr("like_count + 1")).Error
	}
	if err != nil {
		return 0, err
	}
	return int64(like.ID), nil
}

func UnlikeAction(ctx context.Context, targetID int64, likeType int64) error {
	var result *gorm.DB
	switch likeType {
	case VideoLikeType:
		result = DB.WithContext(ctx).Where("target_id = ? AND type = ?", targetID, VideoLikeType).Delete(&Like{})
	case CommentLikeType:
		result = DB.WithContext(ctx).Where("target_id = ? AND type = ?", targetID, CommentLikeType).Delete(&Like{})
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
		err = DB.WithContext(ctx).Model(&Comment{}).Where("id = ?", targetID).Update("like_count", gorm.Expr("like_count - 1")).Error
	}
	if err != nil {
		return err
	}
	return nil
}

func UnlikeActionByID(ctx context.Context, id int64, likeType int64) error {
	var targetID int64
	result := DB.WithContext(ctx).Where("id = ?", id).Pluck("target_id", &targetID)
	if result.Error != nil {
		return result.Error
	}
	result.Delete(&Like{})
	if result.RowsAffected == 0 {
		return errno.LikeIsNotExistErr
	}
	var err error
	switch likeType {
	case VideoLikeType:
		_, err = rpc.UpdateVideoLikeCount(ctx, targetID, -1)
	case CommentLikeType:
		err = DB.WithContext(ctx).Model(&Comment{}).Where("id = ?", targetID).Update("like_count", gorm.Expr("like_count - 1")).Error
	}
	if err != nil {
		return err
	}
	return nil
}

func QueryLikeVideoListByUserID(ctx context.Context, userID int64, pageNum, pageSize int64) ([]int64, error) {
	var likes []int64
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	err := DB.WithContext(ctx).Model(&Like{}).Where("type = 1 AND user_id = ?", userID).Pluck("target_id", &likes).
		Offset(int(pageNum) - 1).Limit(int(pageSize)).Error
	if err != nil {
		return nil, err
	}
	if len(likes) == 0 {
		return nil, errno.LikeIsNotExistErr
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
