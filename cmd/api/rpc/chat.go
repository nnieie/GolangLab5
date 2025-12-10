package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	kitextracing "github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/kitex_gen/chat"
	"github.com/nnieie/golanglab5/kitex_gen/chat/chatservice"
	"github.com/nnieie/golanglab5/pkg/constants"
)

var chatClient chatservice.Client

func InitChatRPC() {
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	c, err := chatservice.NewClient(
		constants.ChatServiceName,
		client.WithRPCTimeout(constants.RPCTimeout),
		client.WithConnectTimeout(constants.ConnectTimeout),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithResolver(r),
		client.WithSuite(kitextracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
	chatClient = c
}

func SendPrivateMessage(ctx context.Context, req *chat.SendPrivateMessageRequest) (*chat.SendPrivateMessageResponse, error) {
	resp, err := chatClient.SendPrivateMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SendGroupMessage(ctx context.Context, req *chat.SendGroupMessageRequest) (*chat.SendGroupMessageResponse, error) {
	resp, err := chatClient.SendGroupMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func QueryPrivateMessageHistory(ctx context.Context, req *chat.QueryPrivateHistoryMessageRequest) (*chat.QueryPrivateHistoryMessageResponse, error) {
	resp, err := chatClient.QueryPrivateHistoryMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func QueryGroupMessageHistory(ctx context.Context, req *chat.QueryGroupHistoryMessageRequest) (*chat.QueryGroupHistoryMessageResponse, error) {
	resp, err := chatClient.QueryGroupHistoryMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func QueryPrivateOfflineMessages(ctx context.Context, req *chat.QueryPrivateOfflineMessageRequest) (*chat.QueryPrivateOfflineMessageResponse, error) {
	resp, err := chatClient.QueryPrivateOfflineMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func QueryGroupOfflineMessages(ctx context.Context, req *chat.QueryGroupOfflineMessageRequest) (*chat.QueryGroupOfflineMessageResponse, error) {
	resp, err := chatClient.QueryGroupOfflineMessage(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func QueryGroupMembers(ctx context.Context, req *chat.QueryGroupMembersRequest) (*chat.QueryGroupMembersResponse, error) {
	resp, err := chatClient.QueryGroupMembers(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func CheckUserExistInGroup(ctx context.Context, req *chat.CheckUserExistInGroupRequest) (*chat.CheckUserExistInGroupResponse, error) {
	resp, err := chatClient.CheckUserExistInGroup(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
