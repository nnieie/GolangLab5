package service

import (
	"github.com/nnieie/golanglab5/internal/video/dal/db"
)

func (s *VideoService) BatchUpdateVideoLikeCount(counts map[int64]int64) error {
	return db.BatchUpdateVideoLikeCount(s.ctx, counts)
}
