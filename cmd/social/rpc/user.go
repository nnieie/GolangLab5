package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/kitex_gen/base"
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
		client.WithResolver(r),
		client.WithRPCTimeout(constants.RPCTimeout),
		client.WithConnectTimeout(constants.ConnectTimeout),
		client.WithFailureRetry(retry.NewFailurePolicy()),
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}

func QueryUsersByIDs(ctx context.Context, userIDs []int64) ([]*base.User, error) {
	resp, err := userClient.QueryUsersByIDs(ctx, &user.QueryUsersByIDsRequest{
		UserIds: userIDs,
	})
	if err != nil {
		return nil, err
	}
	return resp.Users, nil
}
