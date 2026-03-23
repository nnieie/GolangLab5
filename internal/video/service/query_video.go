package service

import (
	"github.com/nnieie/golanglab5/internal/video/dal/db"
	"github.com/nnieie/golanglab5/internal/video/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *VideoService) QueryVideoByID(videoID string) (*base.Video, error) {
	video, err := db.QueryVideoByID(s.ctx, videoID)
	if err != nil {
		return nil, err
	}
	return pack.DBVideoToBaseVideo(video), nil
}

func (s *VideoService) QueryVideosByIDs(videoIDs []string) ([]*base.Video, error) {
	videos, err := db.QueryVideosByIDs(s.ctx, videoIDs)
	if err != nil {
		return nil, err
	}
	return pack.DBVideosToBaseVideos(videos), nil
}
