package main

import (
	"context"

	"github.com/nnieie/golanglab5/cmd/interaction/service"
	interaction "github.com/nnieie/golanglab5/kitex_gen/interaction"
	"github.com/nnieie/golanglab5/pkg/utils"
)

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct{}

// LikeAction implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) LikeAction(ctx context.Context, req *interaction.LikeActionRequest,
) (resp *interaction.LikeActionResponse, err error) {
	resp = new(interaction.LikeActionResponse)
	err = service.NewInteractionService(ctx).LikeAction(req.UserId, req.ActionType, req.VideoId, req.CommentId)
	resp.Base = utils.BuildBaseResp(err)
	return
}

// GetLikeList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetLikeList(ctx context.Context, req *interaction.GetLikeListRequest,
) (resp *interaction.GetLikeListResponse, err error) {
	resp = new(interaction.GetLikeListResponse)
	videos, err := service.NewInteractionService(ctx).GetLikeList(req.UserId, req.PageNum, req.PageSize)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = videos
	return
}

// CommentAction implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CommentAction(ctx context.Context, req *interaction.CommentRequest,
) (resp *interaction.CommentResponse, err error) {
	resp = new(interaction.CommentResponse)
	err = service.NewInteractionService(ctx).PublishComment(req.UserId, req.VideoId, req.CommentId, req.Content)
	resp.Base = utils.BuildBaseResp(err)
	return
}

// GetCommentList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetCommentList(ctx context.Context, req *interaction.GetCommentListRequest,
) (resp *interaction.GetCommentListResponse, err error) {
	resp = new(interaction.GetCommentListResponse)
	comments, err := service.NewInteractionService(ctx).GetCommentList(req.VideoId, req.CommentId, req.PageNum, req.PageSize)
	resp.Base = utils.BuildBaseResp(err)
	resp.Data = comments
	return
}

// DeleteComment implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) DeleteComment(ctx context.Context, req *interaction.DeleteCommentRequest,
) (resp *interaction.DeleteCommentResponse, err error) {
	resp = new(interaction.DeleteCommentResponse)
	err = service.NewInteractionService(ctx).DeleteComment(req.UserId, req.VideoId, req.CommentId)
	resp.Base = utils.BuildBaseResp(err)
	return
}
