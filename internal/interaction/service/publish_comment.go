package service

import (
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/nnieie/golanglab5/internal/interaction/dal/db"
	"github.com/nnieie/golanglab5/pkg/tracer"
)

func (s *interactionService) PublishComment(userID string, videoID, commentID *string, content string) error {
	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}
	if videoID != nil {
		intVideoID, parseErr := strconv.ParseInt(*videoID, 10, 64)
		if parseErr != nil {
			return parseErr
		}
		_, err = db.CreateComment(s.ctx, &db.Comment{
			UserID:  intUserID,
			VideoID: intVideoID,
			Content: content,
		})
		if err == nil {
			tracer.InteractionCommentCounter.Add(s.ctx, 1, metric.WithAttributes(
				attribute.String("action", "add"),
			))
		}
		return err
	} else if commentID != nil {
		intCommentID, parseErr := strconv.ParseInt(*commentID, 10, 64)
		if parseErr != nil {
			return parseErr
		}
		_, err = db.CreateComment(s.ctx, &db.Comment{
			UserID:   intUserID,
			ParentID: intCommentID,
			Content:  content,
		})
		if err == nil {
			tracer.InteractionCommentCounter.Add(s.ctx, 1, metric.WithAttributes(
				attribute.String("action", "add"),
			))
		}
		return err
	}
	return nil
}
