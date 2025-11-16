package service

import (
	"context"
	"time"

	"github.com/nnieie/golanglab5/cmd/interaction/dal/cache"
	"github.com/nnieie/golanglab5/pkg/errno"
	"github.com/nnieie/golanglab5/pkg/kafka"
	"github.com/nnieie/golanglab5/pkg/logger"
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

	// 发送 Kafka 消息
	event := &kafka.LikeEvent{
		BaseEvent: kafka.BaseEvent{
			EventID:   kafka.GenerateEventID(userID, time.Now().UnixMilli()),
			EventType: kafka.EventTypeLike,
			Timestamp: time.Now().UnixMilli(),
			UserID:    userID,
		},
		VideoID:   videoID,
		CommentID: commentID,
		Action:    actionType,
	}

	// 异步发送，不等待结果
	go func() {
		producer := kafka.GetProducer()
		if producer == nil {
			logger.Errorf("Kafka producer not initialized")
			return
		}
		if err := producer.PublishLikeEvent(ctx, event); err != nil {
			logger.Errorf("Failed to publish like event: %v", err)
		}
	}()

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
