package service

import (
	"github.com/nnieie/golanglab5/cmd/user/dal/db"
	"github.com/nnieie/golanglab5/cmd/user/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *UserService) QueryUserByID(userID int64) (*base.User, error) {
	user, err := db.QueryUserByID(s.ctx, userID)
	if err != nil {
		return nil, err
	}
	return pack.DBUserTobaseUser(user), nil
}

func (s *UserService) QueryUsersByIDs(userIDs []int64) ([]*base.User, error) {
	users, err := db.QueryUserByIDList(s.ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return pack.DBUserTobaseUsers(users), nil
}
