package service

import (
	"context"

	"github.com/nnieie/golanglab5/cmd/interaction/dal/cache"
	"github.com/nnieie/golanglab5/cmd/interaction/kafka"
	"github.com/nnieie/golanglab5/pkg/errno"
)

func (s *interactionService) LikeAction(userID int64, actionType int64, videoID, commentID *int64) error {
	ctx := s.ctx
	var exists bool
	var err error
	// 检查重复操作
	switch {
	case videoID != nil:
		exists, err = cache.CheckVideoLikeExists(ctx, userID, *videoID)
	case commentID != nil:
		exists, err = cache.CheckCommentLikeExists(ctx, userID, *commentID)
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
	if videoID != nil {
		if err = s.handleVideoLike(ctx, userID, actionType, *videoID); err != nil {
			return err
		}
	} else if commentID != nil {
		if err = s.handleCommentLike(ctx, userID, actionType, *commentID); err != nil {
			return err
		}
	}

	go kafka.SendLikeEvent(ctx, userID, videoID, commentID, actionType)

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
