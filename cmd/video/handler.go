package main

import (
	"context"

	video "github.com/nnieie/golanglab5/kitex_gen/video"
)

// VideoServiceImpl implements the last service interface defined in the IDL.
type VideoServiceImpl struct{}

// PublishVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) PublishVideo(ctx context.Context, req *video.PublishRequest) (resp *video.PublishResponse, err error) {
	// TODO: Your code here...
	return
}

// GetPublishList implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPublishList(ctx context.Context, req *video.GetPublishListRequest) (resp *video.GetPublishListResponse, err error) {
	// TODO: Your code here...
	return
}

// SearchVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) SearchVideo(ctx context.Context, req *video.SearchVideoRequest) (resp *video.SearchVideoResponse, err error) {
	// TODO: Your code here...
	return
}

// GetPopularVideo implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetPopularVideo(ctx context.Context, req *video.GetPopularVideoListRequest) (resp *video.GetPopularVideoListResponse, err error) {
	// TODO: Your code here...
	return
}

// GetVideoStream implements the VideoServiceImpl interface.
func (s *VideoServiceImpl) GetVideoStream(ctx context.Context, req *video.VideoStreamRequest) (resp *video.VideoStreamResponse, err error) {
	// TODO: Your code here...
	return
}
