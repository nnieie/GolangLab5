package service

import (
	"time"

	"github.com/nnieie/golanglab5/internal/video/dal/db"
	"github.com/nnieie/golanglab5/internal/video/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

const (
	defaultFeedLimit = 20
)

func (s *VideoService) FeedVideo(latestTime *int64) ([]*base.Video, error) {
	if latestTime == nil {
		videos, err := db.QueryLatestVideos(s.ctx, defaultFeedLimit)
		if err != nil {
			return nil, err
		}
		return pack.DBVideosToBaseVideos(videos), nil
	}

	latestTimeTime := time.UnixMilli(*latestTime)

	videos, err := db.QueryVideoByLatestTime(s.ctx, latestTimeTime)
	if err != nil {
		return nil, err
	}
	return pack.DBVideosToBaseVideos(videos), nil
}
