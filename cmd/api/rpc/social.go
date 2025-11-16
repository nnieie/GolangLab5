package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/kitex_gen/social"
	"github.com/nnieie/golanglab5/kitex_gen/social/socialservice"
	"github.com/nnieie/golanglab5/pkg/constants"
)

var socialClient socialservice.Client

func InitSocialRPC() {
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	c, err := socialservice.NewClient(
		constants.SocialServiceName,
		client.WithRPCTimeout(constants.RPCTimeout),
		client.WithConnectTimeout(constants.ConnectTimeout),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithResolver(r),
	)
	if err != nil {
		panic(err)
	}
	socialClient = c
}

func FollowAction(ctx context.Context, req *social.FollowActionRequest) (*social.FollowActionResponse, error) {
	resp, err := socialClient.FollowAction(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetFollowingList(ctx context.Context, req *social.QueryFollowListRequest) (*social.QueryFollowListResponse, error) {
	resp, err := socialClient.QueryFollowList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetFollowerList(ctx context.Context, req *social.QueryFollowerListRequest) (*social.QueryFollowerListResponse, error) {
	resp, err := socialClient.QueryFollowerList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetFriendList(ctx context.Context, req *social.QueryFriendListRequest) (*social.QueryFriendListResponse, error) {
	resp, err := socialClient.QueryFriendList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
