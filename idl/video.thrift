namespace go video

include "base.thrift"

struct PublishRequest{
    1: string title,
    2: string description,
    3: binary video,
    4: binary cover,
    5: i64 user_id,
}

struct PublishResponse{
    1: base.BaseResp base,
}

struct GetPublishListRequest{
    1: string user_id,
    2: i64 page_size,
    3: i64 page_num,
}

struct GetPublishListResponse{
    1:base.BaseResp base,
    2:optional list<base.Video> data,
    3:optional i64 total,
}

struct SearchVideoRequest{
    1: string keyword,
    2: i64 page_num,
    3: i64 page_size,
    4: optional i64 from_date,
    5: optional i64 to_date,
    6: optional string username,
}

struct SearchVideoResponse{
    1:base.BaseResp base,
    2:optional list<base.Video> data,
}

struct GetPopularVideoListRequest{
    1: i64 page_num,
    2: i64 page_size,
}
struct GetPopularVideoListResponse{
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
    PublishResponse PublishVideo(1:PublishRequest req),
    GetPublishListResponse GetPublishList(1:GetPublishListRequest req),
    SearchVideoResponse SearchVideo(1:SearchVideoRequest req),
    GetPopularVideoListResponse GetPopularVideo(1:GetPopularVideoListRequest req),
    VideoStreamResponse GetVideoStream(1:VideoStreamRequest req),
}