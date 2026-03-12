package service

import (
	"strconv"

	"github.com/nnieie/golanglab5/cmd/interaction/dal/db"
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
		_, err := db.CreateComment(s.ctx, &db.Comment{
			UserID:  intUserID,
			VideoID: intVideoID,
			Content: content,
		})
		return err
	} else if commentID != nil {
		intCommentID, parseErr := strconv.ParseInt(*commentID, 10, 64)
		if parseErr != nil {
			return parseErr
		}
		_, err := db.CreateComment(s.ctx, &db.Comment{
			UserID:   intUserID,
			ParentID: intCommentID,
			Content:  content,
		})
		return err
	}
	return nil
}
