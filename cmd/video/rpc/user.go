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
	kitextracing "github.com/kitex-contrib/obs-opentelemetry/tracing"
)

var userClient userservice.Client

func InitUserRPC() {
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
		client.WithSuite(kitextracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
	userClient = c
}

func SearchUserIds(ctx context.Context, username string, pageNum, pageSize int64) ([]int64, error) {
	resp, err := userClient.SearchUserIdsByName(ctx, &user.SearchUserIdsByNameRequest{
		Pattern:  username,
		PageNum:  pageNum,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, err
	}
	return resp.UserIds, nil
}
