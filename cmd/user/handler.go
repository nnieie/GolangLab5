package main

import (
	"bytes"
	"context"

	"github.com/nnieie/golanglab5/cmd/user/service"
	"github.com/nnieie/golanglab5/kitex_gen/base"
	user "github.com/nnieie/golanglab5/kitex_gen/user"
	"github.com/nnieie/golanglab5/pkg/utils"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	Snowflake *utils.Snowflake
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	resp = new(user.RegisterResponse)
	userID, err := service.NewUserService(ctx, s.Snowflake).Register(req.Username, req.Password)
	resp.Base = utils.BuildBaseResp(err)
	resp.UserId = &userID
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	resp = new(user.LoginResponse)
	user, err := service.NewUserService(ctx, s.Snowflake).Login(req.Username, req.Password, req.MFAcode)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = user
	return
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	resp = new(user.UserInfoResponse)
	user, err := service.NewUserService(ctx, s.Snowflake).GetUserInfo(req.UserId)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = user
	return
}

// UploadAvatar implements the UserServiceImpl interface.
func (s *UserServiceImpl) UploadAvatar(ctx context.Context, req *user.UploadAvatarRequest) (resp *user.UploadAvatarResponse, err error) {
	resp = new(user.UploadAvatarResponse)
	avatarData := bytes.NewReader(req.Data)
	user, err := service.NewUserService(ctx, s.Snowflake).UploadAvatar(req.UserId, avatarData, req.FileName)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = user
	return
}

// GetMFAQrcode implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetMFAQrcode(ctx context.Context, req *user.GetMFAQrcodeRequest) (resp *user.GetMFAQrcodeResponse, err error) {
	resp = new(user.GetMFAQrcodeResponse)
	secret, qrcode, err := service.NewUserService(ctx, s.Snowflake).GetMFAqrcode(req.UserId)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = &base.MFAQrcode{
		Secret: secret,
		Qrcode: qrcode,
	}
	return
}

// MFABind implements the UserServiceImpl interface.
func (s *UserServiceImpl) MFABind(ctx context.Context, req *user.MFABindRequest) (resp *user.MFABindResponse, err error) {
	resp = new(user.MFABindResponse)
	_, err = service.NewUserService(ctx, s.Snowflake).BindMFA(req.Code, req.Secret, req.UserId)
	resp.Base = utils.BuildBaseResp(err)
	return
}
