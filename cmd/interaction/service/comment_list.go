package service

import (
	"github.com/nnieie/golanglab5/cmd/interaction/dal/db"
	"github.com/nnieie/golanglab5/cmd/interaction/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
	"github.com/nnieie/golanglab5/pkg/logger"
)

func (s *interactionService) GetCommentList(videoID, commentID *int64, pageNum, pageSize int64) ([]*base.Comment, error) {
	var comments []*db.Comment
	var err error
	if videoID != nil {
		comments, err = db.QueryCommentByVideoID(s.ctx, *videoID, pageNum, pageSize)
	} else if commentID != nil {
		logger.Debugf("GetCommentList by parentID: %d", *commentID)
		comments, err = db.QueryCommentByParentID(s.ctx, *commentID, pageNum, pageSize)
		logger.Debugf("GetCommentList found %d comments", len(comments))
	}
	if err != nil {
		return nil, err
	}
	return pack.DBCommentsToBaseComments(comments), nil
}
