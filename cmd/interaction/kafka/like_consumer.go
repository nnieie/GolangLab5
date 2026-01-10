package kafka

import (
	"context"
	"encoding/json"
	"time"

	"github.com/nnieie/golanglab5/cmd/interaction/dal/db"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
)

// ConsumeLikeEvent 将点赞事件刷新到数据库
func ConsumeLikeEvent() {
	logger.Debugf("Start ConsumeLikeEvent")
	defer logger.Debugf("Exit ConsumeLikeEvent")

	likeCh := KafkaInstance.Consume(context.Background(), constants.LikeTopic, constants.LikeGroupID)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var batch []*LikeEvent

	for {
		select {
		case msg, ok := <-likeCh:
			if !ok {
				logger.Infof("Like event channel closed")
				return
			}
			var event LikeEvent
			if err := json.Unmarshal(msg.V, &event); err != nil {
				logger.Errorf("Failed to unmarshal like event: %v", err)
				continue
			}
			logger.Infof("Received like event: %+v", event)
			batch = append(batch, &event)
		case <-ticker.C:
			if len(batch) > 0 {
				if err := processLikeBatch(context.Background(), batch); err != nil {
					logger.Errorf("Failed to process like batch: %v", err)
				}
				batch = batch[:0]
				logger.Infof("Processed like event batch")
			}
		}
	}
}

// processLikeBatch 处理点赞事件批次
func processLikeBatch(ctx context.Context, batch []*LikeEvent) error {
	var likesToInsert []db.Like
	var likesToDelete []db.Like

	for _, event := range batch {
		var targetID, likeType int64

		switch {
		case event.VideoID != nil:
			targetID = *event.VideoID
			likeType = db.VideoLikeType
		case event.CommentID != nil:
			targetID = *event.CommentID
			likeType = db.CommentLikeType
		default:
			continue
		}

		switch event.Action {
		case 1: // 点赞
			likesToInsert = append(likesToInsert, db.Like{
				UserID:   event.UserID,
				TargetID: targetID,
				Type:     likeType,
			})
		case 2: // 取消点赞
			likesToDelete = append(likesToDelete, db.Like{
				UserID:   event.UserID,
				TargetID: targetID,
				Type:     likeType,
			})
		}
	}

	// 批量插入点赞
	if len(likesToInsert) > 0 {
		if err := db.BatchLikeAction(ctx, likesToInsert); err != nil {
			logger.Errorf("Failed to batch insert likes: %v", err)
			return err
		}
		logger.Infof("Successfully inserted %d likes", len(likesToInsert))
	}

	// 批量删除点赞
	if len(likesToDelete) > 0 {
		if err := db.BatchUnlikeAction(ctx, likesToDelete); err != nil {
			logger.Errorf("Failed to batch delete likes: %v", err)
			return err
		}
		logger.Infof("Successfully deleted %d likes", len(likesToDelete))
	}

	return nil
}
