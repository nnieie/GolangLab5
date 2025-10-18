package rpc

import (
	"context"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/nnieie/golanglab5/kitex_gen/base"
	"github.com/nnieie/golanglab5/kitex_gen/video"
	"github.com/nnieie/golanglab5/kitex_gen/video/videoservice"
	"github.com/nnieie/golanglab5/pkg/constants"
)

var videoClient videoservice.Client

func InitVideoRPC() {
	r, err := etcd.NewEtcdResolver([]string{constants.EtcdAddr})
	if err != nil {
		panic(err)
	}

	c, err := videoservice.NewClient(
		constants.VideoServiceName,
		client.WithResolver(r),
		client.WithRPCTimeout(constants.RPCTimeout),
		client.WithConnectTimeout(constants.ConnectTimeout),
		client.WithFailureRetry(retry.NewFailurePolicy()),
	)
	if err != nil {
		panic(err)
	}
	videoClient = c
}

func QueryVideoByID(ctx context.Context, videoID int64) (*base.Video, error) {
	resp, err := videoClient.QueryVideoByID(ctx, &video.QueryVideoByIDRequest{
		VideoId: videoID,
	})
	if err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, nil
	}
	return resp.Data, nil
}

func QueryVideosByIDs(ctx context.Context, videoIDs []int64) ([]*base.Video, error) {
	resp, err := videoClient.QueryVideosByIDs(ctx, &video.QueryVideosByIDsRequest{
		VideoIds: videoIDs,
	})
	if err != nil {
		return nil, err
	}
	if resp.Data == nil {
		return nil, nil
	}
	return resp.Data, nil
}

func GetVideoLikeCount(ctx context.Context, videoID int64) (int64, error) {
	resp, err := videoClient.GetVideoLikeCount(ctx, &video.GetVideoLikeCountRequest{
		VideoId: videoID,
	})
	if err != nil {
		return 0, err
	}
	if resp.LikeCount == nil {
		return 0, nil
	}
	return *resp.LikeCount, nil
}

func SetVideoLikeCount(ctx context.Context, videoID int64, likeCount int64) error {
	_, err := videoClient.SetVideoLikeCount(ctx, &video.SetVideoLikeCountRequest{
		VideoId:   videoID,
		LikeCount: likeCount,
	})
	return err
}

func UpdateVideoLikeCount(ctx context.Context, videoID int64, delta int64) (*video.UpdateVideoLikeCountResponse, error) {
	resp, err := videoClient.UpdateVideoLikeCount(ctx, &video.UpdateVideoLikeCountRequest{
		VideoId: videoID,
		Delta:   delta,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
