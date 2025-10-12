package service

import (
	"io"

	"github.com/nnieie/golanglab5/cmd/user/dal/db"
	"github.com/nnieie/golanglab5/cmd/user/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *UserService) UploadAvatar(userID int64, avatar io.Reader) (*base.User, error) {
	imgName, err := s.avatarBucket.GenerateImgName()
	if err != nil {
		return nil, err
	}
	fileURL, err := s.avatarBucket.UploadAvatar(imgName, avatar)
	if err != nil {
		return nil, err
	}
	user, err := db.UpdateAvatar(s.ctx, userID, fileURL)
	if err != nil {
		return nil, err
	}
	return pack.DBUserTobaseUser(user), nil
}
