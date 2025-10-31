package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/kitex_gen/user"
	"github.com/nnieie/golanglab5/kitex_gen/user/userservice"
	"github.com/nnieie/golanglab5/pkg/constants"
)

var userClient userservice.Client

func InitUserRPC() {
	config.Init(constants.UserServiceName)
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	c, err := userservice.NewClient(
		constants.UserServiceName,
		client.WithRPCTimeout(constants.RPCTimeout),
		client.WithConnectTimeout(constants.ConnectTimeout),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithResolver(r),
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}

func UserRegister(ctx context.Context, req *user.RegisterRequest) (*user.RegisterResponse, error) {
	resp, err := userClient.Register(ctx, req)

	if err != nil {
		return nil, err
	}
	return resp, nil
}

func UserLogin(ctx context.Context, req *user.LoginRequest) (*user.LoginResponse, error) {
	resp, err := userClient.Login(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetUserInfo(ctx context.Context, req *user.UserInfoRequest) (*user.UserInfoResponse, error) {
	resp, err := userClient.GetUserInfo(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func GetMFAQrcode(ctx context.Context, req *user.GetMFAQrcodeRequest) (*user.GetMFAQrcodeResponse, error) {
	resp, err := userClient.GetMFAQrcode(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func BindMFA(ctx context.Context, req *user.MFABindRequest) (*user.MFABindResponse, error) {
	resp, err := userClient.MFABind(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func UserAvatar(ctx context.Context, req *user.UploadAvatarRequest) (*user.UploadAvatarResponse, error) {
	resp, err := userClient.UploadAvatar(ctx, req)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func UpdateUserLastLogoutTime(ctx context.Context, req *user.UpdateLastLogoutTimeRequest) error {
	_, err := userClient.UpdateLastLogoutTime(ctx, req)
	return err
}
