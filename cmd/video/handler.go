package main

import (
	"bytes"
	"context"

	"github.com/nnieie/golanglab5/cmd/video/service"
	video "github.com/nnieie/golanglab5/kitex_gen/video"
	"github.com/nnieie/golanglab5/pkg/utils"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct {
	Snowflake *utils.Snowflake
}

// PublishVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishVideo(ctx context.Context, req *video.PublishRequest) (resp *video.PublishResponse, err error) {
	resp = new(video.PublishResponse)
	videoData := bytes.NewReader(req.Video)
	err = service.NewVideoService(ctx, s.Snowflake).PublishVideo(req.UserId, videoData, req.FileName, req.Title, req.Description)
	resp.Base = utils.BuildBaseResp(err)
	return
}

// GetPublishList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPublishList(ctx context.Context, req *video.GetPublishListRequest) (resp *video.GetPublishListResponse, err error) {
	resp = new(video.GetPublishListResponse)
	videos, total, err := service.NewVideoService(ctx, s.Snowflake).GetVideoList(req.UserId, req.PageNum, req.PageSize)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = videos
	resp.Total = &total
	return
}

// SearchVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) SearchVideo(ctx context.Context, req *video.SearchVideoRequest) (resp *video.SearchVideoResponse, err error) {
	resp = new(video.SearchVideoResponse)
	videos, total, err := service.NewVideoService(ctx, s.Snowflake).SearchVideo(req.Keywords, req.PageNum, req.PageSize, req.FromDate, req.ToDate, req.Username)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = videos
	resp.Total = &total
	return
}

// GetPopularVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPopularVideo(ctx context.Context, req *video.GetPopularVideoListRequest) (resp *video.GetPopularVideoListResponse, err error) {
	resp = new(video.GetPopularVideoListResponse)
	videos, err := service.NewVideoService(ctx, s.Snowflake).GetPopularVideo(req.PageNum, req.PageSize)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = videos
	return
}

// GetVideoStream implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetVideoStream(ctx context.Context, req *video.VideoStreamRequest) (resp *video.VideoStreamResponse, err error) {
	resp = new(video.VideoStreamResponse)
	videos, err := service.NewVideoService(ctx, s.Snowflake).FeedVideo(req.LatestTime)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = videos
	return
}

func (s *VideoServiceImpl) QueryVideoByID(ctx context.Context, req *video.QueryVideoByIDRequest) (resp *video.QueryVideoByIDResponse, err error) {
	resp = new(video.QueryVideoByIDResponse)
	video, err := service.NewVideoService(ctx, s.Snowflake).QueryVideoByID(req.VideoId)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = video
	return
}

func (s *VideoServiceImpl) QueryVideosByIDs(ctx context.Context, req *video.QueryVideosByIDsRequest) (resp *video.QueryVideosByIDsResponse, err error) {
	resp = new(video.QueryVideosByIDsResponse)
	videos, err := service.NewVideoService(ctx, s.Snowflake).QueryVideosByIDs(req.VideoIds)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = videos
	return
}

// GetVideoLikeCount implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetVideoLikeCount(ctx context.Context, req *video.GetVideoLikeCountRequest) (resp *video.GetVideoLikeCountResponse, err error) {
	resp = new(video.GetVideoLikeCountResponse)
	likeCount, err := service.NewVideoService(ctx, s.Snowflake).GetVideoLikeCount(req.VideoId)
	resp.Base = utils.BuildBaseResp(err)
	resp.LikeCount = &likeCount
	return
}

// SetVideoLikeCount implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) SetVideoLikeCount(ctx context.Context, req *video.SetVideoLikeCountRequest) (resp *video.SetVideoLikeCountResponse, err error) {
	resp = new(video.SetVideoLikeCountResponse)
	err = service.NewVideoService(ctx, s.Snowflake).SetVideoLikeCount(req.VideoId, req.LikeCount)
	resp.Base = utils.BuildBaseResp(err)
	return
}

// UpdateVideoLikeCount implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) UpdateVideoLikeCount(ctx context.Context, req *video.UpdateVideoLikeCountRequest) (
	resp *video.UpdateVideoLikeCountResponse, err error) {
	resp = new(video.UpdateVideoLikeCountResponse)
	err = service.NewVideoService(ctx, s.Snowflake).UpdateVideoLikeCount(req.VideoId, req.Delta)
	resp.Base = utils.BuildBaseResp(err)
	return
}
