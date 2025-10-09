package main

import (
	"context"

	user "github.com/nnieie/golanglab5/kitex_gen/user"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct{}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, req *user.UserInfoRequest) (resp *user.UserInfoResponse, err error) {
	// TODO: Your code here...
	return
}

// UploadAvatar implements the UserServiceImpl interface.
func (s *UserServiceImpl) UploadAvatar(ctx context.Context, req *user.UploadAvatarRequest) (resp *user.UploadAvatarResponse, err error) {
	// TODO: Your code here...
	return
}

// GetMFAQrcode implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetMFAQrcode(ctx context.Context, req *user.GetMFAQrcodeRequest) (resp *user.GetMFAQrcodeResponse, err error) {
	// TODO: Your code here...
	return
}

// MFABind implements the UserServiceImpl interface.
func (s *UserServiceImpl) MFABind(ctx context.Context, req *user.MFABindRequest) (resp *user.MFABindResponse, err error) {
	// TODO: Your code here...
	return
}
