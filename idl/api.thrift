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
    1: string  user_id,
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
    1: string user_id,
    2: i64 page_size,
    3: i64 page_num,
}

struct GetPublishListData{
    1: list<base.Video> items,
    2: optional i64 total,
}

struct GetPublishListResponse{
    1:base.BaseResp base,
    2:GetPublishListData data,
}

struct GetPopularListRequest{
    1: i64 page_size,
    2: i64 page_num,
}

struct GetPopularListData{
    1: list<base.Video> items,
}

struct GetPopularListResponse{
    1:base.BaseResp base,
    2:GetPopularListData data,
}

struct SearchVideoRequest{
    1: i64 page_size,
    2: i64 page_num,
    3: string keywords,
    4: optional i64 from_date,
    5: optional i64 to_date,
    6: optional string username,
}

struct SearchVideoData{
    1: list<base.Video> items,
    2: optional i64 total,
}

struct SearchVideoResponse{
    1:base.BaseResp base,
    2:SearchVideoData data,
}

struct VideoStreamRequest{
    1:optional i64 latest_time,
}

struct VideoStreamData{
    1: list<base.Video> items,
}

struct VideoStreamResponse{
    1:base.BaseResp base,
    2:VideoStreamData data,
}
service VideoService{
    PublishResponse PublishVideo(1:PublishRequest req)(api.post="/video/publish"),
    GetPublishListResponse GetPublishVideoList(1:GetPublishListRequest req)(api.get="/video/list"),
    GetPopularListResponse GetPopularVideo(1:GetPopularListRequest req)(api.get="/video/popular"),
    SearchVideoResponse SearchVideo(1:SearchVideoRequest req)(api.post="/video/search"),
    VideoStreamResponse GetVideoStream(1:VideoStreamRequest req)(api.get="/video/feed")
}


struct LikeActionRequest{
    1: optional string video_id,
    2: optional string comment_id,
    3: i64 action_type,
}
struct LikeActionResponse{
    1: base.BaseResp base,
}

struct GetLikeListRequest{
    1: string user_id,
    2: i64 page_size,
    3: i64 page_num,
}

struct GetLikeListData{
    1: list<base.Video> items,
}

struct GetLikeListResponse{
    1: base.BaseResp base,
    2: GetLikeListData data,
}

struct PublishCommentRequest{
    1: optional string video_id,
    2: optional string comment_id,
    3: string content,
}

struct PublishCommentResponse{
    1: base.BaseResp base,
}

struct GetCommentListRequest{
    1: optional string video_id,
    2: optional string comment_id,
    3: i64 page_size,
    4: i64 page_num,
}

struct GetCommentListData{
    1: list<base.Comment> items,
}

struct GetCommentListResponse{
    1: base.BaseResp base,
    2: GetCommentListData data,
}

struct DeleteCommentRequest{
    1: optional string video_id,
    2: optional string comment_id,
}

struct DeleteCommentResponse{
    1: base.BaseResp base,
}

service InteractionService{
   LikeActionResponse LikeAction(1:LikeActionRequest req)(api.post="/like/action"),
   GetLikeListResponse GetLikeList(1:GetLikeListRequest req)(api.get="/like/list"),
   PublishCommentResponse PublishComment(1:PublishCommentRequest req)(api.post="/comment/publish"),
   GetCommentListResponse GetCommentList(1:GetCommentListRequest req)(api.get="/comment/list"),
   DeleteCommentResponse DeleteComment(1:DeleteCommentRequest req)(api.delete="/comment/delete")
}


struct FollowActionRequest{
    1: string to_user_id,
    2: i64 action_type,
}

struct FollowActionResponse{
    1: base.BaseResp base,
}

struct GetFollowListRequest{  
    1: string user_id,
    2: i64 page_size,
    3: i64 page_num,
}

struct GetFollowListData{
    1: list<base.User> items,
    2: i64 total,
}

struct GetFollowListResponse{
    1: base.BaseResp base,
    2: GetFollowListData data,
}

struct GetFansListRequest{
    1: string user_id,
    2: i64 page_size,
    3: i64 page_num,
}

struct GetFansListData{
    1: list<base.User> items,
    2: i64 total,
}

struct GetFansListResponse{
    1: base.BaseResp base,
    2: GetFansListData data,
}

struct GetFriendListRequest{
    1: i64 page_size,
    2: i64 page_num,
}

struct GetFriendListData{
    1: list<base.User> items,
    2: i64 total,
}

struct GetFriendListResponse{
    1: base.BaseResp base,
    2: GetFriendListData data,
}

service SocialService{
    FollowActionResponse FollowAction(1:FollowActionRequest req)(api.post="/relation/action"),
    GetFollowListResponse GetFollowList(1:GetFollowListRequest req)(api.get="/following/list"),
    GetFansListResponse GetFansList(1:GetFansListRequest req)(api.get="/follower/list"),
    GetFriendListResponse GetFriendList(1:GetFriendListRequest req)(api.get="/friends/list"),
}


struct SendPrivateMessageRequest{
    1: base.PrivateMessage data,
}

struct SendPrivateMessageResponse{
    1: base.BaseResp base,
}

struct QueryPrivateOfflineMessageRequest{
    1: i64 page_num,
    2: i64 page_size,
}

struct QueryPrivateOfflineMessageData{
    1: list<base.PrivateMessage> items,
}

struct QueryPrivateOfflineMessageResponse{
    1: base.BaseResp base,
    2: QueryPrivateOfflineMessageData data,
}

struct QueryPrivateHistoryMessageRequest{
    1: string to_id,
    2: i64 page_num,
    3: i64 page_size,
}

struct QueryPrivateHistoryMessageData{
    1: list<base.PrivateMessage> items,
}

struct QueryPrivateHistoryMessageResponse{
    1: base.BaseResp base,
    2: QueryPrivateHistoryMessageData data,
}

struct SendGroupMessageRequest{
    1: base.GroupMessage data,
}

struct SendGroupMessageResponse{
    1: base.BaseResp base,
}

struct QueryGroupOfflineMessageRequest{
    1: string group_id,
    2: i64 page_num,
    3: i64 page_size,
}

struct QueryGroupOfflineMessageData{
    1: list<base.GroupMessage> items,
}

struct QueryGroupOfflineMessageResponse{
    1: base.BaseResp base,
    2: QueryGroupOfflineMessageData data,
}

struct QueryGroupHistoryMessageRequest{
    1: string group_id,
    2: i64 page_num,
    3: i64 page_size,
}

struct QueryGroupHistoryMessageData{
    1: list<base.GroupMessage> items,
}

struct QueryGroupHistoryMessageResponse{
    1: base.BaseResp base,
    2: QueryGroupHistoryMessageData data,
}

service ChatService{
    base.BaseResp Chat() (api.get="/ws"),
}