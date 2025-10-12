package service

import (
	"github.com/nnieie/golanglab5/cmd/user/dal/db"
	"github.com/nnieie/golanglab5/cmd/user/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *UserService) GetUserInfo(userID int64) (*base.User, error) {
	user, err := db.QueryUserByID(s.ctx, userID)
	if err != nil {
		return nil, err
	}
	return pack.DBUserTobaseUser(user), nil
}
