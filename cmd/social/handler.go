package main

import (
	"context"

	social "github.com/nnieie/golanglab5/kitex_gen/social"
)

// SocialServiceImpl implements the last service interface defined in the IDL.
type SocialServiceImpl struct{}

// FollowAction implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) FollowAction(ctx context.Context, req *social.FollowActionRequest) (resp *social.FollowActionResponse, err error) {
	// TODO: Your code here...
	return
}

// QueryFollowList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) QueryFollowList(ctx context.Context, req *social.QueryFollowListRequest) (resp *social.QueryFollowListResponse, err error) {
	// TODO: Your code here...
	return
}

// QueryFollowerList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) QueryFollowerList(ctx context.Context, req *social.QueryFollowerListRequest) (resp *social.QueryFollowerListResponse, err error) {
	// TODO: Your code here...
	return
}

// QueryFriendList implements the SocialServiceImpl interface.
func (s *SocialServiceImpl) QueryFriendList(ctx context.Context, req *social.QueryFriendListRequest) (resp *social.QueryFriendListResponse, err error) {
	// TODO: Your code here...
	return
}
