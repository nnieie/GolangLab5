package service

import "context"

type interactionService struct {
	ctx context.Context
}

func NewInteractionService(ctx context.Context) *interactionService {
	return &interactionService{
		ctx: ctx,
	}
}
