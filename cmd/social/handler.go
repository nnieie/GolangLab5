package main

import (
	"context"

	"github.com/nnieie/golanglab5/cmd/social/service"
	social "github.com/nnieie/golanglab5/kitex_gen/social"
	"github.com/nnieie/golanglab5/pkg/utils"
)

// SocialServiceImpl implements the last service interface defined in the IDL.
type SocialServiceImpl struct{}

// FollowAction implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) FollowAction(ctx context.Context, req *social.FollowActionRequest) (resp *social.FollowActionResponse, err error) {
	resp = new(social.FollowActionResponse)
	err = service.NewSocialService(ctx).FollowAction(req.UserId, req.ToUserId, req.ActionType)
	resp.Base = utils.BuildBaseResp(err)
	return
}

// QueryFollowList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) QueryFollowList(ctx context.Context, req *social.QueryFollowListRequest) (resp *social.QueryFollowListResponse, err error) {
	resp = new(social.QueryFollowListResponse)
	following, total, err := service.NewSocialService(ctx).GetFollowingList(req.UserId, req.PageNum, req.PageSize)
	if err != nil {
		resp.Base = utils.BuildBaseResp(err)
		return
	}
	resp.Data = following
	resp.Total = &total
	return
}

// QueryFollowerList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) QueryFollowerList(ctx context.Context, req *social.QueryFollowerListRequest) (resp *social.QueryFollowerListResponse, err error) {
	resp = new(social.QueryFollowerListResponse)
	followers, total, err := service.NewSocialService(ctx).GetFollowerList(req.UserId, req.PageNum, req.PageSize)
	if err != nil {
		resp.Base = utils.BuildBaseResp(err)
		return
	}
	resp.Data = followers
	resp.Total = &total
	return
}

// QueryFriendList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) QueryFriendList(ctx context.Context, req *social.QueryFriendListRequest) (resp *social.QueryFriendListResponse, err error) {
	resp = new(social.QueryFriendListResponse)
	friends, total, err := service.NewSocialService(ctx).GetFriendList(req.UserId, req.PageNum, req.PageSize)
	if err != nil {
		resp.Base = utils.BuildBaseResp(err)
		return
	}
	resp.Data = friends
	resp.Total = &total
	return
}
