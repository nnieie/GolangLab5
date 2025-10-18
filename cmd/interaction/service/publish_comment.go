package service

import "github.com/nnieie/golanglab5/cmd/interaction/dal/db"

func (s *interactionService) PublishComment(userID int64, videoID, commentID *int64, content string) error {
	if videoID != nil {
		_, err := db.CreateComment(s.ctx, &db.Comment{
			UserID:  userID,
			VideoID: *videoID,
			Content: content,
		})
		return err
	} else if commentID != nil {
		_, err := db.CreateComment(s.ctx, &db.Comment{
			UserID:   userID,
			ParentID: *commentID,
			Content:  content,
		})
		return err
	}
	return nil
}
