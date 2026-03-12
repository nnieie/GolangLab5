package rpc

import (
	"context"
	"strconv"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/retry"
	kitextracing "github.com/kitex-contrib/obs-opentelemetry/tracing"
	etcd "github.com/kitex-contrib/registry-etcd"

	"github.com/nnieie/golanglab5/config"
	"github.com/nnieie/golanglab5/kitex_gen/base"
	"github.com/nnieie/golanglab5/kitex_gen/video"
	"github.com/nnieie/golanglab5/kitex_gen/video/videoservice"
	"github.com/nnieie/golanglab5/pkg/constants"
)

var videoClient videoservice.Client

func InitVideoRPC() {
	r, err := etcd.NewEtcdResolver([]string{config.Etcd.Addr})
	if err != nil {
		panic(err)
	}

	c, err := videoservice.NewClient(
		constants.VideoServiceName,
		client.WithResolver(r),
		client.WithRPCTimeout(constants.RPCTimeout),
		client.WithConnectTimeout(constants.ConnectTimeout),
		client.WithFailureRetry(retry.NewFailurePolicy()),
		client.WithSuite(kitextracing.NewClientSuite()),
	)
	if err != nil {
		panic(err)
	}
	videoClient = c
}

func QueryVideoByID(ctx context.Context, videoID int64) (*base.Video, error) {
	resp, err := videoClient.QueryVideoByID(ctx, &video.QueryVideoByIDRequest{
		VideoId: strconv.FormatInt(videoID, 10),
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
	videoIDsStr := make([]string, 0, len(videoIDs))
	for _, videoID := range videoIDs {
		videoIDsStr = append(videoIDsStr, strconv.FormatInt(videoID, 10))
	}
	resp, err := videoClient.QueryVideosByIDs(ctx, &video.QueryVideosByIDsRequest{
		VideoIds: videoIDsStr,
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
		VideoId: strconv.FormatInt(videoID, 10),
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
		VideoId:   strconv.FormatInt(videoID, 10),
		LikeCount: likeCount,
	})
	return err
}

func UpdateVideoLikeCount(ctx context.Context, videoID int64, delta int64) (*video.UpdateVideoLikeCountResponse, error) {
	resp, err := videoClient.UpdateVideoLikeCount(ctx, &video.UpdateVideoLikeCountRequest{
		VideoId: strconv.FormatInt(videoID, 10),
		Delta:   delta,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func BatchUpdateVideoLikeCount(ctx context.Context, videoLikeCounts map[int64]int64) (*video.BatchUpdateVideoLikeCountResponse, error) {
	resp, err := videoClient.BatchUpdateVideoLikeCount(ctx, &video.BatchUpdateVideoLikeCountRequest{
		VideoLikeCounts: videoLikeCounts,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}
