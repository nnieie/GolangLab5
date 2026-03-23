package service

import (
	"github.com/nnieie/golanglab5/internal/user/dal/db"
	"github.com/nnieie/golanglab5/internal/user/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *UserService) GetUserInfo(userID string) (*base.User, error) {
	user, err := db.QueryUserByID(s.ctx, userID)
	if err != nil {
		return nil, err
	}
	return pack.DBUserTobaseUser(user), nil
}
