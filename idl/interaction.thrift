namespace go interaction

include "base.thrift"

struct LikeActionRequest{
    1: optional string video_id,
    2: optional string comment_id,
    3: i64 action_type,
    4: string user_id
}

struct LikeActionResponse{
    1: base.BaseResp base,
}

struct GetLikeListRequest{
    1: string user_id,
    2: i64 page_num,
    3: i64 page_size,
    4: optional i64 last_id,
}

struct GetLikeListData{
    1: list<base.Video> items,
}

struct GetLikeListResponse{
    1: base.BaseResp base,
    2: GetLikeListData data,
}

struct CommentRequest{
    1: string user_id,
    2: optional string video_id,
    3: optional string comment_id,
    4: string content,
}

struct CommentResponse{
    1: base.BaseResp base,
    2: optional string comment_id,
}

struct GetCommentListRequest{
    1: optional string video_id,
    2: optional string comment_id,
    3: i64 page_num,
    4: i64 page_size,
    5: optional i64 last_id,
}

struct GetCommentListData{
    1: list<base.Comment> items,
}

struct GetCommentListResponse{
    1: base.BaseResp base,
    2: GetCommentListData data,
}

struct DeleteCommentRequest{
    1: string user_id,
    2: optional string video_id,
    3: optional string comment_id,
}

struct DeleteCommentResponse{
    1: base.BaseResp base,
}

service InteractionService{
   LikeActionResponse LikeAction(1:LikeActionRequest req),
   GetLikeListResponse GetLikeList(1:GetLikeListRequest req),
   CommentResponse CommentAction(1:CommentRequest req),
   GetCommentListResponse GetCommentList(1:GetCommentListRequest req),
   DeleteCommentResponse DeleteComment(1:DeleteCommentRequest req),
}