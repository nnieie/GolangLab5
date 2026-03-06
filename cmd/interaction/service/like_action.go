package service

import (
	"context"
	"strconv"

	"github.com/nnieie/golanglab5/cmd/interaction/dal/cache"
	"github.com/nnieie/golanglab5/cmd/interaction/kafka"
	"github.com/nnieie/golanglab5/pkg/errno"
)

func (s *interactionService) LikeAction(userID string, actionType int64, videoID, commentID *string) error {
	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}
	var intVideoID *int64
	if videoID != nil {
		parsedVideoID, parseErr := strconv.ParseInt(*videoID, 10, 64)
		if parseErr != nil {
			return parseErr
		}
		intVideoID = &parsedVideoID
	}
	var intCommentID *int64
	if commentID != nil {
		parsedCommentID, parseErr := strconv.ParseInt(*commentID, 10, 64)
		if parseErr != nil {
			return parseErr
		}
		intCommentID = &parsedCommentID
	}
	ctx := s.ctx
	var exists bool

	// 检查重复操作
	switch {
	case intVideoID != nil:
		exists, err = cache.CheckVideoLikeExists(ctx, intUserID, *intVideoID)
	case intCommentID != nil:
		exists, err = cache.CheckCommentLikeExists(ctx, intUserID, *intCommentID)
	default:
		return nil
	}
	if err != nil {
		return err
	}
	if actionType == 1 && exists {
		return errno.LikeAlreadyExistErr
	}
	if actionType == 2 && !exists {
		return errno.LikeIsNotExistErr
	}

	//  更新 Redis
	if intVideoID != nil {
		if err = s.handleVideoLike(ctx, intUserID, actionType, *intVideoID); err != nil {
			return err
		}
	} else if intCommentID != nil {
		if err = s.handleCommentLike(ctx, intUserID, actionType, *intCommentID); err != nil {
			return err
		}
	}

	go kafka.SendLikeEvent(context.Background(), intUserID, intVideoID, intCommentID, actionType)

	return nil
}

// handleVideoLike 处理视频点赞逻辑
func (s *interactionService) handleVideoLike(ctx context.Context, userID, actionType int64, videoID int64) error {
	delta := int64(1)
	if actionType == 2 {
		delta = -1
	}

	// 更新点赞数
	if err := cache.IncrVideoLikeCount(ctx, videoID, delta); err != nil {
		return err
	}

	// 记录用户点赞关系
	if actionType == 1 {
		if err := cache.SetVideoLike(ctx, userID, videoID); err != nil {
			return err
		}
	} else {
		if err := cache.DelVideoLike(ctx, userID, videoID); err != nil {
			return err
		}
	}
	return nil
}

// handleCommentLike 处理评论点赞逻辑
func (s *interactionService) handleCommentLike(ctx context.Context, userID, actionType int64, commentID int64) error {
	delta := int64(1)
	if actionType == 2 {
		delta = -1
	}

	// 更新点赞数
	if err := cache.IncrCommentLikeCount(ctx, commentID, delta); err != nil {
		return err
	}

	// 记录用户点赞关系
	if actionType == 1 {
		if err := cache.SetCommentLike(ctx, userID, commentID); err != nil {
			return err
		}
	} else {
		if err := cache.DelCommentLike(ctx, userID, commentID); err != nil {
			return err
		}
	}
	return nil
}
