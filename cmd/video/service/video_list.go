package service

import (
	"strconv"

	"github.com/nnieie/golanglab5/cmd/video/dal/db"
	"github.com/nnieie/golanglab5/cmd/video/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func (s *VideoService) GetVideoList(userID string, pageNum, pageSize int64) ([]*base.Video, int64, error) {
	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return nil, 0, err
	}
	videos, err := db.QueryVideoByUserID(s.ctx, intUserID, pageNum, pageSize)
	if err != nil {
		return nil, 0, err
	}
	count, err := db.QueryVideoCountByUserID(s.ctx, intUserID)
	if err != nil {
		return nil, 0, err
	}
	return pack.DBVideosToBaseVideos(videos), count, nil
}
