namespace go social

include "base.thrift"

struct FollowActionRequest{
    1: string to_user_id,
    2: i64 action_type,
    3: string user_id,
}

struct FollowActionResponse{
    1: base.BaseResp base,
}

struct QueryFollowListRequest{
    1: string user_id,
    2: i64 page_num,
    3: i64 page_size,
    4: optional i64 last_id,
}

struct QueryFollowListData{
    1: list<base.User> items,
    2: optional i64 total,
}

struct QueryFollowListResponse{
    1: base.BaseResp base,
    2: QueryFollowListData data,
}

struct QueryFollowerListRequest{
    1: string user_id,
    2: i64 page_num,
    3: i64 page_size,
    4: optional i64 last_id,
}

struct QueryFollowerListData{
    1: list<base.User> items,
    2: optional i64 total,
}

struct QueryFollowerListResponse{
    1: base.BaseResp base,
    2: QueryFollowerListData data,
}

struct QueryFriendListRequest{
    1: i64 page_num,
    2: i64 page_size,
    3: string user_id,
    4: optional i64 last_id,
}

struct QueryFriendListData{
    1: list<base.User> items,
    2: optional i64 total,
}

struct QueryFriendListResponse{
    1: base.BaseResp base,
    2: QueryFriendListData data,
}

service SocialService{
    FollowActionResponse FollowAction(1:FollowActionRequest req),
    QueryFollowListResponse QueryFollowList(1:QueryFollowListRequest req),
    QueryFollowerListResponse QueryFollowerList(1:QueryFollowerListRequest req),
    QueryFriendListResponse QueryFriendList(1:QueryFriendListRequest req),
}
