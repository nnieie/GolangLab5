namespace go video

include "base.thrift"

struct PublishRequest{
    1: string title,
    2: string description,
    3: binary video,
    4: binary cover,
    5: string user_id,
    6: string file_name,
}

struct PublishResponse{
    1: base.BaseResp base,
}

struct GetPublishListRequest{
    1: string user_id,
    2: i64 page_size,
    3: i64 page_num,
    4: optional i64 last_id,
}

struct GetPublishListData{
    1: list<base.Video> items,
    2: optional i64 total,
}

struct GetPublishListResponse{
    1:base.BaseResp base,
    2:GetPublishListData data,
}

struct SearchVideoRequest{
    1: string keywords,
    2: i64 page_num,
    3: i64 page_size,
    4: optional i64 from_date,
    5: optional i64 to_date,
    6: optional string username,
    7: optional i64 last_id,
}

struct SearchVideoData{
    1: list<base.Video> items,
    2: optional i64 total,
}

struct SearchVideoResponse{
    1:base.BaseResp base,
    2:SearchVideoData data,
}

struct GetPopularVideoListRequest{
    1: i64 page_num,
    2: i64 page_size,
    3: optional i64 last_id,
}

struct GetPopularVideoListData{
    1: list<base.Video> items,
}

struct GetPopularVideoListResponse{
    1:base.BaseResp base,
    2:GetPopularVideoListData data,
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

struct QueryVideoByIDRequest{
    1: string video_id,
}

struct QueryVideoByIDResponse{
    1: base.BaseResp base,
    2: optional base.Video data,
}

struct QueryVideosByIDsRequest{
    1: list<string> video_ids,
}

struct QueryVideosByIDsResponse{
    1: base.BaseResp base,
    2: optional list<base.Video> data,
}

struct GetVideoLikeCountRequest{
    1: string video_id,
}

struct GetVideoLikeCountResponse{
    1: base.BaseResp base,
    2: optional i64 like_count,
}

struct SetVideoLikeCountRequest{
    1: string video_id,
    2: i64 like_count,
}

struct SetVideoLikeCountResponse{
    1: base.BaseResp base,
}

struct UpdateVideoLikeCountRequest{
    1: string video_id,
    2: i64 delta,
}

struct UpdateVideoLikeCountResponse{
    1: base.BaseResp base,
}

struct BatchUpdateVideoLikeCountRequest{
    1: map<i64, i64> video_like_counts,
}

struct BatchUpdateVideoLikeCountResponse{
    1: base.BaseResp base,
}

service VideoService{
    PublishResponse PublishVideo(1:PublishRequest req),
    GetPublishListResponse GetPublishList(1:GetPublishListRequest req),
    SearchVideoResponse SearchVideo(1:SearchVideoRequest req),
    GetPopularVideoListResponse GetPopularVideo(1:GetPopularVideoListRequest req),
    VideoStreamResponse GetVideoStream(1:VideoStreamRequest req),
    QueryVideoByIDResponse QueryVideoByID(1:QueryVideoByIDRequest req),
    QueryVideosByIDsResponse QueryVideosByIDs(1:QueryVideosByIDsRequest req),
    GetVideoLikeCountResponse GetVideoLikeCount(1:GetVideoLikeCountRequest req),
    SetVideoLikeCountResponse SetVideoLikeCount(1:SetVideoLikeCountRequest req),
    UpdateVideoLikeCountResponse UpdateVideoLikeCount(1:UpdateVideoLikeCountRequest req),
    BatchUpdateVideoLikeCountResponse BatchUpdateVideoLikeCount(1:BatchUpdateVideoLikeCountRequest req),
}