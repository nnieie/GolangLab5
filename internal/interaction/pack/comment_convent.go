package pack

import (
	"strconv"

	"github.com/nnieie/golanglab5/internal/interaction/dal/db"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func DBCommentToBaseComment(comment *db.Comment) *base.Comment {
	if comment == nil {
		return nil
	}
	return &base.Comment{
		Id:         strconv.FormatUint(uint64(comment.ID), 10),
		UserId:     strconv.FormatInt(comment.UserID, 10),
		VideoId:    strconv.FormatInt(comment.VideoID, 10),
		ParentId:   strconv.FormatInt(comment.ParentID, 10),
		Content:    comment.Content,
		LikeCount:  comment.LikeCount,
		ChildCount: comment.ChildCount,
		CreatedAt:  comment.CreatedAt.String(),
		UpdatedAt:  comment.UpdatedAt.String(),
	}
}

func DBCommentsToBaseComments(comments []*db.Comment) []*base.Comment {
	baseComments := make([]*base.Comment, 0, len(comments))
	for _, comment := range comments {
		baseComments = append(baseComments, DBCommentToBaseComment(comment))
	}
	return baseComments
}
