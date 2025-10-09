namespace go interaction

include "base.thrift"

struct LikeActionRequest{
    1: optional i64 video_id,
    2: optional i64 comment_id,
    3: i64 action_type,
    4: i64 user_id
}

struct LikeActionResponse{
    1: base.BaseResp base,
}

struct GetLikeListRequest{
    1: i64 user_id,
    2: i64 page_num,
    3: i64 page_size,
}

struct GetLikeListResponse{
    1: base.BaseResp base,
    2: optional list<base.Video> data,
}

struct CommentRequest{
    1: i64 user_id,
    2: optional i64 video_id,
    3: optional i64 comment_id,
    4: string content,
}

struct CommentResponse{
    1: base.BaseResp base,
    2: optional i64 comment_id,
}

struct GetCommentListRequest{
    1: optional i64 video_id,
    2: optional i64 comment_id,
    3: i64 page_num,
    4: i64 page_size,
}

struct GetCommentListResponse{
    1: base.BaseResp base,
    2: optional list<base.Comment> data,
}

struct DeleteCommentRequest{
    1: i64 user_id,
    2: optional i64 video_id,
    3: optional i64 comment_id,
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