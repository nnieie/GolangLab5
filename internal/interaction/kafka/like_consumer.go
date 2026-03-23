package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/nnieie/golanglab5/internal/interaction/dal/db"
	"github.com/nnieie/golanglab5/pkg/constants"
	"github.com/nnieie/golanglab5/pkg/logger"
	"github.com/nnieie/golanglab5/pkg/tracer"
)

const maxBatchSize = 500

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

			if len(batch) >= maxBatchSize {
				if err := processLikeBatch(context.Background(), batch); err != nil {
					logger.Errorf("Failed to process like batch (size triggered): %v", err)
				} else {
					logger.Infof("Processed like event batch (size triggered): %d records", len(batch))
				}
				batch = batch[:0]
			}
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
	// 使用 map 进行状态折叠
	type actionState struct {
		like   db.Like
		action int64
	}
	latestState := make(map[string]actionState)

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

		key := fmt.Sprintf("%d_%d_%d", event.UserID, targetID, likeType)
		latestState[key] = actionState{
			like: db.Like{
				UserID:   event.UserID,
				TargetID: targetID,
				Type:     likeType,
			},
			action: event.Action,
		}
	}

	var likesToInsert []db.Like
	var likesToDelete []db.Like

	for _, state := range latestState {
		switch state.action {
		case 1:
			likesToInsert = append(likesToInsert, state.like)
		case 2:
			likesToDelete = append(likesToDelete, state.like)
		}
	}

	// 批量插入点赞
	if len(likesToInsert) > 0 {
		if err := db.BatchLikeAction(ctx, likesToInsert); err != nil {
			logger.Errorf("Failed to batch insert likes: %v", err)
			tracer.MQConsumeCounter.Add(ctx, 1, metric.WithAttributes(
				attribute.String("topic", constants.LikeTopic),
				attribute.String("status", "fail"),
				attribute.String("error_type", "db_insert_error"),
			))
			return err
		}
		logger.Infof("Successfully inserted %d likes", len(likesToInsert))
		tracer.MQConsumeCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("topic", constants.LikeTopic),
			attribute.String("status", "success"),
		))
		tracer.InteractionLikeCounter.Add(ctx, int64(len(likesToInsert)), metric.WithAttributes(
			attribute.String("action", "like"),
		))
	}

	// 批量删除点赞
	if len(likesToDelete) > 0 {
		if err := db.BatchUnlikeAction(ctx, likesToDelete); err != nil {
			logger.Errorf("Failed to batch delete likes: %v", err)
			tracer.MQConsumeCounter.Add(ctx, 1, metric.WithAttributes(
				attribute.String("topic", constants.LikeTopic),
				attribute.String("status", "fail"),
				attribute.String("error_type", "db_delete_error"),
			))
			return err
		}
		logger.Infof("Successfully deleted %d likes", len(likesToDelete))
		tracer.MQConsumeCounter.Add(ctx, 1, metric.WithAttributes(
			attribute.String("topic", constants.LikeTopic),
			attribute.String("status", "success"),
		))
		tracer.InteractionLikeCounter.Add(ctx, int64(len(likesToDelete)), metric.WithAttributes(
			attribute.String("action", "unlike"),
		))
	}

	return nil
}
