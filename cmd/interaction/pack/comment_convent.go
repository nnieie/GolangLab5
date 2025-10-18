package pack

import (
	"github.com/nnieie/golanglab5/cmd/interaction/dal/db"
	"github.com/nnieie/golanglab5/kitex_gen/base"
)

func DBCommentToBaseComment(comment *db.Comment) *base.Comment {
	if comment == nil {
		return nil
	}
	return &base.Comment{
		Id:         int64(comment.ID),
		UserId:     comment.UserID,
		VideoId:    comment.VideoID,
		ParentId:   comment.ParentID,
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
