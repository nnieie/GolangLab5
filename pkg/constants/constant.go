package constants

import "time"

const (
	UserTableName    = "users"
	VideoTableName   = "videos"
	CommentTableName = "comments"
	LikeTableName    = "likes"
	FollowsTableName = "follows"

	PrivateMessageTableName = "private_messages"
	GroupMessageTableName   = "group_messages"
	LastLogoutTimeTableName = "last_logout_times"
	GroupMembersTableName   = "group_members"
)

const (
	DBMaxOpenConns    = 100
	DBMaxIdleConns    = 20
	DBConnMaxLifetime = time.Hour
	DBConnMaxIdleTime = 10 * time.Minute
)

const (
	IdentityKey           = "user_id"
	AccessTokenSecretKey  = "nnieie"
	RefreshTokenSecretKey = "nnieiee"
)

const (
	MaxAvatarSize = 5 << 20
	MaxVideoSize  = 4 << 30
)

const (
	EtcdAddr       = "127.0.0.1:2379"
	RPCTimeout     = 3 * time.Second
	ConnectTimeout = 50 * time.Millisecond
)

const (
	APIServiceName         = "api"
	UserServiceName        = "user"
	VideoServiceName       = "video"
	SocialServiceName      = "social"
	InteractionServiceName = "interaction"
	ChatServiceName        = "chat"
)

const (
	AvatarBucketName = "tkrpc-avatar"
	VideoBucketName  = "tkrpc-video"

	AvatarPublicDomain = "avatar.tkrpc.miaooo.qzz.io"
	VideoPublicDomain  = "video.tkrpc.miaooo.qzz.io"
)

const (
	TOTPSecret        = "totp_secret:"
	TOTPSecretExpTime = time.Minute * 15
)

const (
	CryptKey = "aes256114514aes2561919810aes2566" // 32 bytes key for AES-256
)

const (
	WriteWait      = 10 * time.Second
	PongWait       = 60 * time.Second
	PingPeriod     = 50 * time.Second
	MaxMessageSize = 1024
)
