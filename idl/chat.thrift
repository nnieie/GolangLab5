namespace go chat

include "base.thrift"

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
    2: optional list<base.Message> data,
}

struct QueryPrivateHistoryMessageRequest{
    1:i64 user_id,
    2:i64 to_id,
    3:i64 page_num,
    4:i64 page_size,
}

struct QueryPrivateHistoryMessageResponse{
    1:base.BaseResp base,
    2:optional list<base.Message> data,
}

struct QueryGroupHistoryMessageRequest{
    1:i64 user_id,
    2:i64 to_id,
    3:i64 page_num,
    4:i64 page_size,
}

struct QueryGroupHistoryMessageResponse{
    1:base.BaseResp base,
    2:optional list<base.Message> data,
}

service ChatService{
    SendMessageResponse SendMessage (1:SendMessageRequest req),
    QueryPrivateOfflineMessageResponse QueryPrivateOfflineMessage (1:QueryPrivateOfflineMessageRequest req),
    QueryPrivateHistoryMessageResponse QueryPrivateHistoryMessage(1:QueryPrivateHistoryMessageRequest req),
    QueryGroupHistoryMessageResponse QueryGroupHistoryMessage(1:QueryGroupHistoryMessageRequest req),
}
