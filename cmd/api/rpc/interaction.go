package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/nnieie/golanglab5/kitex_gen/interaction"
	"github.com/nnieie/golanglab5/kitex_gen/interaction/interactionservice"
	"github.com/nnieie/golanglab5/pkg/constants"
)

var interactionClient interactionservice.Client

func InitInteractionRPC() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddr})
	if err != nil {
		panic(err)
	}

	c, err := interactionservice.NewClient(
		constants.InteractionServiceName,
		client.WithRPCTimeout(constants.RPCTimeout),
		client.WithConnectTimeout(constants.ConnectTimeout),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithResolver(r),
	)
	if err != nil {
		panic(err)
	}
	interactionClient = c
}

func LikeAction(ctx context.Context, req *interaction.LikeActionRequest) (*interaction.LikeActionResponse, error) {
	resp, err := interactionClient.LikeAction(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetLikeList(ctx context.Context, req *interaction.GetLikeListRequest) (*interaction.GetLikeListResponse, error) {
	resp, err := interactionClient.GetLikeList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func CommentAction(ctx context.Context, req *interaction.CommentRequest) (*interaction.CommentResponse, error) {
	resp, err := interactionClient.CommentAction(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetCommentList(ctx context.Context, req *interaction.GetCommentListRequest) (*interaction.GetCommentListResponse, error) {
	resp, err := interactionClient.GetCommentList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func DeleteComment(ctx context.Context, req *interaction.DeleteCommentRequest) (*interaction.DeleteCommentResponse, error) {
	resp, err := interactionClient.DeleteComment(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
