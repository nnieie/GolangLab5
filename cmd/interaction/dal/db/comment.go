package db

import (
	"context"

	"gorm.io/gorm"

	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/logger"
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

func CreateComment(c context.Context, comment *Comment) (int64, error) {
	err := DB.WithContext(c).Create(comment).Error
	if err != nil {
		return 0, err
	}

	if comment.ParentID != 0 {
		DB.WithContext(c).Where("id = ?", comment.ParentID).Update("child_count", gorm.Expr("child_count + 1"))
	}
	return int64(comment.ID), nil
}

func QueryCommentByID(c context.Context, id int64) (*Comment, error) {
	var comment Comment
	err := DB.WithContext(c).Model(&Comment{}).Where("id = ?", id).Find(&comment).Error
	if err != nil {
		return nil, err
	}
	if comment == (Comment{}) {
		return nil, errno.CommentIsNotExistErr
	}
	return &comment, nil
}

func QueryCommentByVideoID(c context.Context, videoID int64, pageNum, pageSize int64) ([]*Comment, error) {
	var comments []*Comment
	err := DB.WithContext(c).Model(&Comment{}).Where("video_id = ?", videoID).Limit(int(pageSize)).Offset(int(pageNum) - 1).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, errno.CommentIsNotExistErr
	}
	return comments, nil
}

func QueryCommentByParentID(c context.Context, parentID int64, pageNum, pageSize int64) ([]*Comment, error) {
	var comments []*Comment
	err := DB.WithContext(c).Model(&Comment{}).Where("parent_id = ?", parentID).Limit(int(pageSize)).Offset(int(pageNum) - 1).Find(&comments).Error
	if err != nil {
		return nil, err
	}
	if len(comments) == 0 {
		return nil, errno.CommentIsNotExistErr
	}
	return comments, nil
}

func DeleteCommentByCommentID(c context.Context, userID, commentID int64) error {
	result := DB.WithContext(c).Where("user_id = ? AND id = ?", userID, commentID).Delete(&Comment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errno.CommentIsNotExistErr
	}
	result = DB.WithContext(c).Where("parent_id = ?", commentID).Delete(&Comment{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteCommentByVideoID(c context.Context, userID, videoID int64) error {
	var comments []*Comment
	result := DB.WithContext(c).Where("user_id = ? AND video_id = ?", userID, videoID).Find(&comments).Delete(&Comment{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errno.CommentIsNotExistErr
	}
	for _, comment := range comments {
		logger.Debugf("Deleting child comments of comment ID: %d", comment.ID)
		if err := DB.WithContext(c).Where("parent_id = ?", comment.ID).Delete(&Comment{}).Error; err != nil {
			return err
		}
	}
	return nil
}
