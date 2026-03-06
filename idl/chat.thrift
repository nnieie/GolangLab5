namespace go chat

include "base.thrift"

struct SendPrivateMessageRequest{
    1: base.PrivateMessage data,
}

struct SendPrivateMessageResponse{
    1: base.BaseResp base,
}

struct QueryPrivateOfflineMessageRequest{
    1: string user_id,
    2: string to_user_id,
    3: i64 page_num,
    4: i64 page_size,
}

struct QueryPrivateOfflineMessageData{
    1: list<base.PrivateMessage> items,
}

struct QueryPrivateOfflineMessageResponse{
    1: base.BaseResp base,
    2: QueryPrivateOfflineMessageData data,
}

struct QueryPrivateHistoryMessageRequest{
    1:string user_id,
    2:string to_user_id,
    3:i64 page_num,
    4:i64 page_size,
}

struct QueryPrivateHistoryMessageData{
    1: list<base.PrivateMessage> items,
}

struct QueryPrivateHistoryMessageResponse{
    1:base.BaseResp base,
    2:QueryPrivateHistoryMessageData data,
}

struct QueryGroupHistoryMessageRequest{
    1:string user_id,
    2:string group_id,
    3:i64 page_num,
    4:i64 page_size,
}

struct SendGroupMessageRequest{
    1: base.GroupMessage data,
}

struct SendGroupMessageResponse{
    1: base.BaseResp base,
}

struct QueryGroupOfflineMessageRequest{
    1: string user_id,
    2: string group_id,
    3: i64 page_num,
    4: i64 page_size,
}

struct QueryGroupOfflineMessageData{
    1: list<base.GroupMessage> items,
}

struct QueryGroupOfflineMessageResponse{
    1: base.BaseResp base,
    2: QueryGroupOfflineMessageData data,
}

struct QueryGroupHistoryMessageData{
    1: list<base.GroupMessage> items,
}

struct QueryGroupHistoryMessageResponse{
    1:base.BaseResp base,
    2:QueryGroupHistoryMessageData data,
}

struct QueryGroupMembersRequest{
    1: string group_id,
}

struct QueryGroupMembersResponse{
    1: base.BaseResp base,
    2: optional list<i64> members,
}

struct CheckUserExistInGroupRequest{
    1: string user_id,
    2: string group_id,
}

struct CheckUserExistInGroupResponse{
    1: base.BaseResp base,
    2: optional bool exist,
}

service ChatService{
    SendPrivateMessageResponse SendPrivateMessage (1:SendPrivateMessageRequest req),
    QueryPrivateOfflineMessageResponse QueryPrivateOfflineMessage (1:QueryPrivateOfflineMessageRequest req),
    QueryPrivateHistoryMessageResponse QueryPrivateHistoryMessage(1:QueryPrivateHistoryMessageRequest req),
    SendGroupMessageResponse SendGroupMessage (1:SendGroupMessageRequest req),
    QueryGroupOfflineMessageResponse QueryGroupOfflineMessage (1:QueryGroupOfflineMessageRequest req),
    QueryGroupHistoryMessageResponse QueryGroupHistoryMessage(1:QueryGroupHistoryMessageRequest req),
    QueryGroupMembersResponse QueryGroupMembers(1:QueryGroupMembersRequest req),
    CheckUserExistInGroupResponse CheckUserExistInGroup(1:CheckUserExistInGroupRequest req),
}
