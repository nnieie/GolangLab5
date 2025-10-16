package service

import (
	"io"
	"strings"

	"github.com/nnieie/golanglab5/cmd/video/dal/db"
)

// TODO: 分片上传
func (s *VideoService) PublishVideo(userID int64, video io.Reader, fileName string, title, description string) error {
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
		UserID:      userID,
		VideoURL:    fileURL,
		Title:       title,
		Description: description,
	})
	if err != nil {
		return err
	}
	return nil
}
