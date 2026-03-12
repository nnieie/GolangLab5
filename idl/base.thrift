namespace go base

struct BaseResp {
    1: i64 code,
    2: string msg,
}

struct User{
    1: string id,
    2: string username,
    3: string avatar_url,
    4: string created_at,
    5: string updated_at,
    6: string deleted_at,
}

struct Video{
  1: string id,
    2: string user_id,
    3: string video_url,
    4: string cover_url,
    5: string title,
    6: string description,
    7: i64 visit_count,
    8: i64 like_count,
    9: i64 comment_count,
    10: string created_at,
    11: string updated_at,
    12: string deleted_at
}

struct Comment{
  1: string id,
    2: string user_id,
  3: string video_id,
  4: string parent_id,
    5: i64 like_count,
    6: i64 child_count,
    7: string content,
    8: string created_at,
    9: string updated_at,
    10: string deleted_at,
}

struct PrivateMessage {
  1: string id,
    2: string from_user_id,
    3: string to_user_id,
  4: string content,
  5: i64 created_at,
}

struct GroupMessage {
  1: string id,
    2: string from_user_id,
  3: string group_id,
  4: string content,
  5: i64 created_at,
}

struct MFAQrcode{
    1: string secret,
    2: string qrcode,
}