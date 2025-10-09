package main

import (
	"context"

	interaction "github.com/nnieie/golanglab5/kitex_gen/interaction"
)

// InteractionServiceImpl implements the last service interface defined in the IDL.
type InteractionServiceImpl struct{}

// LikeAction implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) LikeAction(ctx context.Context, req *interaction.LikeActionRequest,
) (resp *interaction.LikeActionResponse, err error) {
	// TODO: Your code here...
	return
}

// GetLikeList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetLikeList(ctx context.Context, req *interaction.GetLikeListRequest,
) (resp *interaction.GetLikeListResponse, err error) {
	// TODO: Your code here...
	return
}

// CommentAction implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) CommentAction(ctx context.Context, req *interaction.CommentRequest,
) (resp *interaction.CommentResponse, err error) {
	// TODO: Your code here...
	return
}

// GetCommentList implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) GetCommentList(ctx context.Context, req *interaction.GetCommentListRequest,
) (resp *interaction.GetCommentListResponse, err error) {
	// TODO: Your code here...
	return
}

// DeleteComment implements the InteractionServiceImpl interface.
func (s *InteractionServiceImpl) DeleteComment(ctx context.Context, req *interaction.DeleteCommentRequest,
) (resp *interaction.DeleteCommentResponse, err error) {
	// TODO: Your code here...
	return
}
