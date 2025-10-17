package service

import "context"

type socialService struct {
	ctx context.Context
}

func NewSocialService(ctx context.Context) *socialService {
	return &socialService{
		ctx: ctx,
	}
}
