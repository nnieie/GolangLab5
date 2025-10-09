namespace go api

include "base.thrift"

struct RegisterRequest{
    1: string username,
    2: string password,
}

struct RegisterResponse{
    1: base.BaseResp base,
}

struct LoginRequest{
    1: string username,
    2: string password,
    3: optional string code,
}

struct LoginResponse{
    1: base.BaseResp base,
    2: base.User data,
}

struct GetUserInfoRequest{
    1: i64  user_id,
}
struct GetUserInfoResponse{
    1: base.BaseResp base,
    2: base.User data,
}

struct UploadAvatarRequest{
    1: binary data,
}
struct UploadAvatarResponse{
    1: base.BaseResp base,
    2: base.User data,
}

struct GetMFAResponse{
    1: base.BaseResp base,
    2: base.MFAQrcode data,
}

struct MFABindRequest{
    1: string code,
    2: string secret,
}

struct MFABindResponse{
    1: base.BaseResp base,
}

service UserService {
    RegisterResponse Register (1: RegisterRequest req)(api.post="/user/register"),
    LoginResponse Login(1: LoginRequest req)(api.post="/user/login"),
    GetUserInfoResponse GetUserInfo(1:GetUserInfoRequest req)(api.get="/user/info"),
    UploadAvatarResponse UploadAvatar(1:UploadAvatarRequest req)(api.put="/user/avatar/upload"),
    GetMFAResponse GetMFA()(api.get="/auth/mfa/qrcode"),
    MFABindResponse MFABind(1:MFABindRequest req)(api.post="/auth/mfa/bind")
}


struct PublishRequest{
    1: string title,
    2: string description,
    3: binary data,
}

struct PublishResponse{
    1: base.BaseResp base,
}

struct GetPublishListRequest{
    1: i64 user_id,
    2: i64 page_size,
    3: i64 page_num,
}

struct GetPublishListResponse{
    1:base.BaseResp base,
    2:optional list<base.Video> data,
}

struct GetPopularListRequest{
    1: i64 page_size,
    2: i64 page_num,
}
struct GetPopularListResponse{
    1:base.BaseResp base,
    2:optional list<base.Video> data,
}

struct SearchVideoRequest{
    1: i64 page_size,
    2: i64 page_num,
    3: string keyword,
    4: optional i64 from_date,
    5: optional i64 to_date,
    6: optional string username,
}

struct SearchVideoResponse{
    1:base.BaseResp base,
    2:optional list<base.Video> data,
}

struct VideoStreamRequest{
    1:optional i64 latest_time,
}

struct VideoStreamResponse{
    1:base.BaseResp base,
    2:optional list<base.Video> data,
}
service VideoService{
    PublishResponse PublishVideo(1:PublishRequest req)(api.post="/video/publish"),
    GetPublishListResponse GetPublishVideoList(1:GetPublishListRequest req)(api.get="/video/list"),
    GetPopularListResponse GetPopularVideo(1:GetPopularListRequest req)(api.get="/video/popular"),
    SearchVideoResponse SearchVideo(1:SearchVideoRequest req)(api.post="/video/search"),
    VideoStreamResponse GetVideoStream(1:VideoStreamRequest req)(api.get="/video/feed")
}


struct LikeActionRequest{
    1: optional i64 video_id,
    2: optional i64 comment_id,
    3: i64 action_type,
}
struct LikeActionResponse{
    1: base.BaseResp base,
}

struct GetLikeListRequest{
    1: i64 user_id,
    2: i64 page_size,
    3: i64 page_num,
}

struct GetLikeListResponse{
    1: base.BaseResp base,
    2: list<base.Video> data,
}

struct CommentRequest{
    1: optional i64 video_id,
    2: optional i64 comment_id,
    3: string content,
}

struct CommentResponse{
    1: base.BaseResp base,
}

struct GetCommentListRequest{
    1: optional i64 video_id,
    2: optional i64 comment_id,
    3: i64 page_size,
    4: i64 page_num,
}

struct GetCommentListResponse{
    1: base.BaseResp base,
    2: list<base.Comment> data,
}

struct DeleteCommentRequest{
    1: optional i64 video_id,
    2: optional i64 comment_id,
}

struct DeleteCommentResponse{
    1: base.BaseResp base,
}

service InteractionService{
   LikeActionResponse LikeAction(1:LikeActionRequest req)(api.post="/like/action"),
   GetLikeListResponse GetLikeList(1:GetLikeListRequest req)(api.get="/like/list"),
   CommentResponse CommentVideo(1:CommentRequest req)(api.post="/comment/publish"),
   GetCommentListResponse GetCommentList(1:GetCommentListRequest req)(api.get="/comment/list"),
   DeleteCommentResponse DeleteComment(1:DeleteCommentRequest req)(api.delete="/comment/delete")
}


struct FollowActionRequest{
    1: i64 to_user_id,
    2: i64 action_type,
}

struct FollowActionResponse{
    1: base.BaseResp base,
}

struct GetFollowListRequest{  
    1: i64 user_id,
    2: i64 page_size,
    3: i64 page_num,
}

struct GetFollowListResponse{
    1: base.BaseResp base,
    2: list<base.User> data,
}

struct GetFansListRequest{
    1: i64 user_id,
    2: i64 page_size,
    3: i64 page_num,
}

struct GetFansListResponse{
    1: base.BaseResp base,
    2: list<base.User> data,
}

struct GetFriendListRequest{
    1: i64 page_size,
    2: i64 page_num,
}

struct GetFriendListResponse{
    1: base.BaseResp base,
    2: list<base.User> data,
}

service SocialService{
    FollowActionResponse FollowAction(1:FollowActionRequest req)(api.post="/relation/action"),
    GetFollowListResponse GetFollowList(1:GetFollowListRequest req)(api.get="/following/list"),
    GetFansListResponse GetFansList(1:GetFansListRequest req)(api.get="/follower/list"),
    GetFriendListResponse GetFriendList(1:GetFriendListRequest req)(api.get="/friends/list"),
}


struct SendMessageRequest{
    1: i64 id,
    2: i64 user_id,
    3: i64 to_user_id,
    4: i64 group_id,
    5: i64 type,
    6: string content,
    7: i64 create_time,
}

struct SendMessageResponse{
    1: base.BaseResp base,
}

struct QueryPrivateOfflineMessageRequest{
    1: i64 user_id,
}

struct QueryPrivateOfflineMessageResponse{
    1: base.BaseResp base,
    2: list<base.Message> data,
}

struct QueryPrivateHistoryMessageRequest{
    1: i64 user_id,
    2: i64 to_id,
    3: i64 page_num,
    4: i64 page_size,
}

struct QueryPrivateHistoryMessageResponse{
    1: base.BaseResp base,
    2: list<base.Message> data,
}

struct QueryGroupHistoryMessageRequest{
    1: i64 user_id,
    2: i64 to_id,
    3: i64 page_num,
    4: i64 page_size,
}

struct QueryGroupHistoryMessageResponse{
    1: base.BaseResp base,
    2: list<base.Message> data,
}

service ChatService{
    base.BaseResp Chat() (api.get="/ws"),
}