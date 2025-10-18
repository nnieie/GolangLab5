package service

import "github.com/nnieie/golanglab5/cmd/interaction/dal/db"

func (s *interactionService) DeleteComment(userID int64, videoID, commentID *int64) error {
	if videoID != nil {
		return db.DeleteCommentByVideoID(s.ctx, userID, *videoID)
	} else if commentID != nil {
		return db.DeleteCommentByCommentID(s.ctx, userID, *commentID)
	}
	return nil
}
