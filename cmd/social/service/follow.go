package service

import "github.com/nnieie/golanglab5/cmd/social/dal/db"

func (s *socialService) FollowAction(userID, toUserID, actionType int64) error {
	switch actionType {
	case 0:
		return db.CreateFollows(s.ctx, &db.Follow{
			UserID:     toUserID,
			FollowerID: userID,
		})
	case 1:
		return db.DeleteFollows(s.ctx, &db.Follow{
			UserID:     toUserID,
			FollowerID: userID,
		})
	default:
		return nil
	}
}
