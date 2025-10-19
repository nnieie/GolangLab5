package service

import (
	"time"

	"github.com/nnieie/golanglab5/cmd/user/dal/db"
)

func (s *UserService) GetLastLogoutTime(userID int64) (int64, error) {
	logoutTime, err := db.QueryLastLogoutTime(s.ctx, userID)
	return logoutTime.Unix(), err
}

func (s *UserService) UpdateLastLogoutTime(userID int64, logoutTime int64) error {
	return db.UpdateLastLogoutTime(s.ctx, userID, time.Unix(logoutTime, 0))
}
