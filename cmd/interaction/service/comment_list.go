package service

import (
	"strconv"

	"github.com/nnieie/golanglab5/cmd/interaction/dal/db"
	"github.com/nnieie/golanglab5/cmd/interaction/pack"
	"github.com/nnieie/golanglab5/kitex_gen/base"
	"github.com/nnieie/golanglab5/pkg/logger"
)

func (s *interactionService) GetCommentList(videoID, commentID *string, pageNum, pageSize int64) ([]*base.Comment, error) {
	var comments []*db.Comment
	var err error
	if videoID != nil {
		intVideoID, parseErr := strconv.ParseInt(*videoID, 10, 64)
		if parseErr != nil {
			return nil, parseErr
		}
		comments, err = db.QueryCommentByVideoID(s.ctx, intVideoID, pageNum, pageSize)
	} else if commentID != nil {
		intCommentID, parseErr := strconv.ParseInt(*commentID, 10, 64)
		if parseErr != nil {
			return nil, parseErr
		}
		logger.Debugf("GetCommentList by parentID: %d", intCommentID)
		comments, err = db.QueryCommentByParentID(s.ctx, intCommentID, pageNum, pageSize)
		logger.Debugf("GetCommentList found %d comments", len(comments))
	}
	if err != nil {
		return nil, err
	}
	return pack.DBCommentsToBaseComments(comments), nil
}
