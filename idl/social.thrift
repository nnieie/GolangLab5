namespace go social

include "base.thrift"

struct FollowActionRequest{
    1: i64 to_user_id,
    2: i64 action_type,
    3: i64 user_id,
}

struct FollowActionResponse{
    1: base.BaseResp base,
}

struct QueryFollowListRequest{
    1: i64 user_id,
    2: i64 page_num,
    3: i64 page_size,
}

struct QueryFollowListResponse{
    1: base.BaseResp base,
    2: optional list<base.User> data,
    3: optional i64 total,
}

struct QueryFollowerListRequest{
    1: i64 user_id,
    2: i64 page_num,
    3: i64 page_size,
}

struct QueryFollowerListResponse{
    1: base.BaseResp base,
    2: optional list<base.User> data,
    3: optional i64 total,
}

struct QueryFriendListRequest{
    1: i64 page_num,
    2: i64 page_size,
    3: i64 user_id,
}

struct QueryFriendListResponse{
    1: base.BaseResp base,
    2: optional list<base.User> data,
    3: optional i64 total,
}

service SocialService{
    FollowActionResponse FollowAction(1:FollowActionRequest req),
    QueryFollowListResponse QueryFollowList(1:QueryFollowListRequest req),
    QueryFollowerListResponse QueryFollowerList(1:QueryFollowerListRequest req),
    QueryFriendListResponse QueryFriendList(1:QueryFriendListRequest req),
}
