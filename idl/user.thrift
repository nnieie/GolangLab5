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
    1:string avatar_url,
    2:i64 user_id,
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

service UserService {
    RegisterResponse Register (1: RegisterRequest req),
    LoginResponse Login(1: LoginRequest req),
    UserInfoResponse GetUserInfo(1:UserInfoRequest req),
    UploadAvatarResponse UploadAvatar(1:UploadAvatarRequest req),
    GetMFAQrcodeResponse GetMFAQrcode(1:GetMFAQrcodeRequest req),
    MFABindResponse MFABind(1:MFABindRequest req),
}
