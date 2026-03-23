package service

import (
	"io"
	"strconv"

	"github.com/nnieie/golanglab5/internal/video/dal/db"
	"github.com/nnieie/golanglab5/pkg/tracer"
)

// TODO: 分片上传
func (s *VideoService) PublishVideo(userID string, video io.Reader, fileName, videoURL string, title, description string) error {
	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}

	_, err = db.CreateVideo(s.ctx, &db.Video{
		UserID:      intUserID,
		VideoURL:    videoURL,
		Title:       title,
		Description: description,
	})
	if err != nil {
		return err
	}
	tracer.VideoPublishCounter.Add(s.ctx, 1)
	return nil
}
