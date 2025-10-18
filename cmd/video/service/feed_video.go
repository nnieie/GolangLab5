package service

import (
	"time"

	"github.com/nnieie/golanglab5/cmd/video/dal/db"
	"github.com/nnieie/golanglab5/cmd/video/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *VideoService) FeedVideo(latestTime *int64) ([]*base.Video, error) {
	if latestTime == nil {
		return nil, nil
	}
	latestTimeTime := time.Unix(*latestTime, 0)
	videos, err := db.QueryVideoByLatestTime(s.ctx, latestTimeTime)
	if err != nil {
		return nil, err
	}
	return pack.DBVideosToBaseVideos(videos), nil
}
