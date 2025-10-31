package db

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/errno"
)

type Comment struct {
	UserID     int64
	VideoID    int64
	ParentID   int64
	LikeCount  int64
	ChildCount int64
	Content    string
	gorm.Model
}

func (Comment) TableName() string {
	return constants.CommentTableName
}

func CreateComment(ctx context.Context, comment *Comment) (int64, error) {
	// 把所有操作都放进一个事务里
	err := DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 在事务中创建评论
		if err := tx.Create(comment).Error; err != nil {
			return err
		}

		// 在事务中更新父评论计数
		if comment.ParentID != 0 {
			if err := tx.Model(&Comment{}).Where("id = ?", comment.ParentID).Update("child_count", gorm.Expr("child_count + 1")).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return int64(comment.ID), nil
}

func QueryCommentByID(ctx context.Context, id int64) (*Comment, error) {
	var comment Comment
	err := DB.WithContext(ctx).Model(&Comment{}).Where("id = ?", id).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.CommentIsNotExistErr
		}
		return nil, err
	}
	return &comment, nil
}

func QueryCommentByVideoID(ctx context.Context, videoID int64, pageNum, pageSize int64) ([]*Comment, error) {
	var comments []*Comment
	err := DB.WithContext(ctx).Model(&Comment{}).Where("video_id = ?", videoID).Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func QueryCommentByParentID(ctx context.Context, parentID int64, pageNum, pageSize int64) ([]*Comment, error) {
	var comments []*Comment
	err := DB.WithContext(ctx).Model(&Comment{}).Where("parent_id = ?", parentID).Limit(int(pageSize)).Offset(int((pageNum - 1) * pageSize)).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func DeleteCommentByCommentID(ctx context.Context, userID, commentID int64) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var commentToDelete Comment
		// 找到要删除的评论
		if err := tx.Where("id = ? AND user_id = ?", commentID, userID).First(&commentToDelete).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errno.CommentIsNotExistErr
			}
			return err
		}

		// 删除这条评论
		if err := tx.Delete(&commentToDelete).Error; err != nil {
			return err
		}

		// 如果它有父评论，就给父评论的 child_count - 1
		if commentToDelete.ParentID != 0 {
			if err := tx.Model(&Comment{}).Where("id = ?", commentToDelete.ParentID).Update("child_count", gorm.Expr("child_count - 1")).Error; err != nil {
				return err
			}
		}

		// 4. 删除该评论的所有子评论
		if err := tx.Where("parent_id = ?", commentID).Delete(&Comment{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func DeleteCommentsByVideoID(ctx context.Context, userID, videoID int64) error {
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 找到该 user_id 发表的评论
		var commentIDs []int64
		err := tx.Model(&Comment{}).Where("user_id = ? AND video_id = ?", userID, videoID).Pluck("id", &commentIDs).Error
		if err != nil {
			return err
		}

		// 如果该用户在这个视频下没有任何评论，直接成功返回
		if len(commentIDs) == 0 {
			return nil
		}

		// 找出这些评论下的所有子评论ID
		var childCommentIDs []int64
		if err := tx.Model(&Comment{}).
			Where("parent_id IN ?", commentIDs).
			Pluck("id", &childCommentIDs).Error; err != nil {
			return err
		}

		// 更新父评论的 child_count
		type ParentUpdateInfo struct {
			ParentID int64
			Count    int64
		}
		var updates []ParentUpdateInfo
		err = tx.Model(&Comment{}).
			Select("parent_id, COUNT(*) as count").
			Where("id IN ? AND parent_id != 0", commentIDs).
			Group("parent_id").
			Find(&updates).Error
		if err != nil {
			return err
		}
		for _, update := range updates {
			if err := tx.Model(&Comment{}).Where("id = ?", update.ParentID).Update("child_count", gorm.Expr("child_count - ?", update.Count)).Error; err != nil {
				return err
			}
		}

		// 批量删除
		commentIDs = append(commentIDs, childCommentIDs...)
		if err := tx.Where("id IN ?", commentIDs).Delete(&Comment{}).Error; err != nil {
			return err
		}

		return nil
	})
}
