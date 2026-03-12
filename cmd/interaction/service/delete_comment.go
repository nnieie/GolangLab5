package service

import (
	"strconv"

	"github.com/nnieie/golanglab5/cmd/interaction/dal/db"
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
		return db.DeleteCommentsByVideoID(s.ctx, intUserID, intVideoID)
	} else if commentID != nil {
		intCommentID, parseErr := strconv.ParseInt(*commentID, 10, 64)
		if parseErr != nil {
			return parseErr
		}
		return db.DeleteCommentByCommentID(s.ctx, intUserID, intCommentID)
	}
	return nil
}
