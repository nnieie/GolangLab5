package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/kitex_gen/video"
	"github.com/nnieie/golanglab5/kitex_gen/video/videoservice"
	"github.com/nnieie/golanglab5/pkg/constants"
)

var videoClient videoservice.Client

func InitVideoRPC() {
	config.Init(constants.VideoServiceName)
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	c, err := videoservice.NewClient(
		constants.VideoServiceName,
		client.WithRPCTimeout(constants.RPCTimeout),
		client.WithConnectTimeout(constants.ConnectTimeout),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithResolver(r),
	)
	if err != nil {
		panic(err)
	}
	videoClient = c
}

func GetPopularVideo(ctx context.Context, req *video.GetPopularVideoListRequest) (*video.GetPopularVideoListResponse, error) {
	resp, err := videoClient.GetPopularVideo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetVideoStream(ctx context.Context, req *video.VideoStreamRequest) (*video.VideoStreamResponse, error) {
	resp, err := videoClient.GetVideoStream(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func GetPublishList(ctx context.Context, req *video.GetPublishListRequest) (*video.GetPublishListResponse, error) {
	resp, err := videoClient.GetPublishList(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func SearchVideo(ctx context.Context, req *video.SearchVideoRequest) (*video.SearchVideoResponse, error) {
	resp, err := videoClient.SearchVideo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func PublishVideo(ctx context.Context, req *video.PublishRequest) (*video.PublishResponse, error) {
	resp, err := videoClient.PublishVideo(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
