package service

import "github.com/nnieie/golanglab5/cmd/video/dal/db"

func (s *VideoService) GetVideoLikeCount(videoID int64) (int64, error) {
	return db.QueryVideoLikeCount(s.ctx, videoID)
}

func (s *VideoService) SetVideoLikeCount(videoID, count int64) error {
	return db.SetVideoLikeCount(s.ctx, videoID, count)
}

func (s *VideoService) UpdateVideoLikeCount(videoID, delta int64) error {
	return db.UpdateVideoLikeCount(s.ctx, videoID, delta)
}
