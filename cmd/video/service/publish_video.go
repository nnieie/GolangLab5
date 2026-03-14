package service

import (
	"io"
	"strconv"
	"strings"

	"github.com/nnieie/golanglab5/cmd/video/dal/db"
	"github.com/nnieie/golanglab5/pkg/tracer"
)

// TODO: 分片上传
func (s *VideoService) PublishVideo(userID string, video io.Reader, fileName string, title, description string) error {
	intUserID, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		return err
	}
	videoName, err := s.videoBucket.GenerateVideoName()
	if err != nil {
		return err
	}
	videoName = strings.Join([]string{videoName, fileName}, "_")
	fileURL, err := s.videoBucket.UploadVideo(videoName, video)
	if err != nil {
		return err
	}
	_, err = db.CreateVideo(s.ctx, &db.Video{
		UserID:      intUserID,
		VideoURL:    fileURL,
		Title:       title,
		Description: description,
	})
	if err != nil {
		return err
	}
	tracer.VideoPublishCounter.Add(s.ctx, 1)
	return nil
}
