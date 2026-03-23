package service

import (
	"github.com/nnieie/golanglab5/internal/user/dal/db"
	"github.com/nnieie/golanglab5/internal/user/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *UserService) UploadAvatar(userID string, avatarurl string) (*base.User, error) {
	user, err := db.UpdateAvatar(s.ctx, userID, avatarurl)
	if err != nil {
		return nil, err
	}
	return pack.DBUserTobaseUser(user), nil
}
