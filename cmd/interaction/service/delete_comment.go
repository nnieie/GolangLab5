package service

import (
	"strconv"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/nnieie/golanglab5/cmd/interaction/dal/db"
	"github.com/nnieie/golanglab5/pkg/tracer"
)

func (s *interactionService) DeleteComment(userID string, videoID, commentID *string) error {
	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}
	if videoID != nil {
		intVideoID, parseErr := strconv.ParseInt(*videoID, 10, 64)
		if parseErr != nil {
			return parseErr
		}
		err = db.DeleteCommentsByVideoID(s.ctx, intUserID, intVideoID)
		if err == nil {
			tracer.InteractionCommentCounter.Add(s.ctx, 1, metric.WithAttributes(
				attribute.String("action", "delete"),
			))
		}
		return err
	} else if commentID != nil {
		intCommentID, parseErr := strconv.ParseInt(*commentID, 10, 64)
		if parseErr != nil {
			return parseErr
		}
		err = db.DeleteCommentByCommentID(s.ctx, intUserID, intCommentID)
		if err == nil {
			tracer.InteractionCommentCounter.Add(s.ctx, 1, metric.WithAttributes(
				attribute.String("action", "delete"),
			))
		}
		return err
	}
	return nil
}
