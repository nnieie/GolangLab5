package service

import "github.com/nnieie/golanglab5/cmd/interaction/dal/db"

func (s *interactionService) LikeAction(userID, actionType int64, videoID, commentID *int64) error {
	var err error
	var targetID int64
	var likeType int64
	switch {
	case videoID != nil:
		targetID = *videoID
		likeType = db.VideoLikeType
	case commentID != nil:
		targetID = *commentID
		likeType = db.CommentLikeType
	default:
		return nil
	}

	switch actionType {
	case db.LikeActionType:
		_, err = db.LikeAction(s.ctx, &db.Like{
			UserID:   userID,
			TargetID: targetID,
			Type:     likeType,
		})
	case db.UnlikeActionType:
		err = db.UnlikeAction(s.ctx, targetID, likeType)
	}

	return err
}
