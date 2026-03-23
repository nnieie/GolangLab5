package service

import (
	"strconv"

	"github.com/nnieie/golanglab5/internal/social/dal/db"
)

func (s *socialService) FollowAction(userID, toUserID string, actionType int64) error {
	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}
	intToUserID, err := strconv.ParseInt(toUserID, 10, 64)
	if err != nil {
		return err
	}
	switch actionType {
	case 0:
		return db.CreateFollows(s.ctx, &db.Follow{
			UserID:     intToUserID,
			FollowerID: intUserID,
		})
	case 1:
		return db.DeleteFollows(s.ctx, &db.Follow{
			UserID:     intToUserID,
			FollowerID: intUserID,
		})
	default:
		return nil
	}
}
