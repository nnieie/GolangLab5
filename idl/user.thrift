namespace go user

include "base.thrift"

struct RegisterRequest{
    1: string username,
    2: string password,
}

struct RegisterResponse{
    1: base.BaseResp base,
    2: optional i64 user_id,
}

struct LoginRequest{
    1: string username,
    2: string password,
    3: optional string MFAcode,
}

struct LoginResponse{
    1: base.BaseResp base,
    2: optional base.User data,
}

struct UserInfoRequest{
    1:i64 user_id,
}

struct UserInfoResponse{
    1: base.BaseResp base,
    2: optional base.User data,
}

struct UploadAvatarRequest{
    1:binary data,
    2:i64 user_id,
    3:string file_name,
}

struct UploadAvatarResponse{
    1: base.BaseResp base,
    2: optional base.User data,
}

struct GetMFAQrcodeRequest{
    1: i64 user_id,
}

struct GetMFAQrcodeResponse{
    1: base.BaseResp base,
    2: optional base.MFAQrcode data,
}

struct MFABindRequest{
    1: string code,
    2: string secret,
    3: i64 user_id,
}

struct MFABindResponse{
    1: base.BaseResp base,
}

struct SearchUserIdsByNameRequest {
    1: string pattern,
    2: i64 page_num,
    3: i64 page_size,
}

struct SearchUserIdsByNameResponse {
    1: base.BaseResp base,
    2: optional list<i64> user_ids,
}

struct QueryUserByIDRequest {
    1: i64 user_id,
}

struct QueryUserByIDResponse {
    1: base.BaseResp base,
    2: optional base.User user,
}

struct QueryUsersByIDsRequest {
    1: list<i64> user_ids,
}

struct QueryUsersByIDsResponse {
    1: base.BaseResp base,
    2: optional list<base.User> users,
}

service UserService {
    RegisterResponse Register (1: RegisterRequest req),
    LoginResponse Login(1: LoginRequest req),
    UserInfoResponse GetUserInfo(1:UserInfoRequest req),
    UploadAvatarResponse UploadAvatar(1:UploadAvatarRequest req),
    GetMFAQrcodeResponse GetMFAQrcode(1:GetMFAQrcodeRequest req),
    MFABindResponse MFABind(1:MFABindRequest req),
    SearchUserIdsByNameResponse SearchUserIdsByName(1: SearchUserIdsByNameRequest req),
    QueryUserByIDResponse QueryUserByID(1: QueryUserByIDRequest req),
    QueryUsersByIDsResponse QueryUsersByIDs(1: QueryUsersByIDsRequest req),
}
