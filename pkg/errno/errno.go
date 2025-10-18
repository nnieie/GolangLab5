package errno

import (
	"errors"
	"fmt"
)

const (
	SuccessCode    = 10000
	ServiceErrCode = iota + 10000
	ParamErrCode
	AuthorizationFailedErrCode
	InvalidFileTypeErrCode

	UserAlreadyExistErrCode
	UserIsNotExistErrCode
	UserIsNotExistOrPasswordErrCode

	FollowRelationAlreadyExistErrCode
	FollowRelationNotExistErrCode

	FavoriteRelationAlreadyExistErrCode
	FavoriteRelationNotExistErrCode
	FavoriteActionErrCode

	MessageAddFailedErrCode
	FriendIsNotExistErrCode

	VideoIsNotExistErrCode
	CommentIsNotExistErrCode
	LikeIsNotExistErrCode
	LikeAlreadyExistErrCode
	FollowIsNotExistErrCode
	FollowAlreadyExistErrCode
	ContentIsEmptyErrCode

	NotFriendErrCode
	NotInGroupErrCode

	TOTPSecretExpiredErrCode
	NotGenerateTotpErrCode
	MFAInvalidCodeErrCode
)

const (
	SuccessMsg                     = "success"
	ServerErrMsg                   = "service is unable to start successfully"
	ParamErrMsg                    = "wrong parameter has been given"
	InvalidFileTypeErrMsg          = "invalid file type"
	UserIsNotExistErrMsg           = "user is not exist"
	UserIsNotExistOrPasswordErrMsg = "username or password error"
	FavoriteActionErrMsg           = "favorite add failed"

	MessageAddFailedErrMsg   = "message add failed"
	FriendIsNotExistMsg      = "friend is not exist"
	VideoIsNotExistErrMsg    = "video is not exist"
	CommentIsNotExistErrMsg  = "comment is not exist"
	LikeIsNotExistErrMsg     = "like is not exist"
	LikeAlreadyExistErrMsg   = "like already exists"
	FollowIsNotExistErrMsg   = "follow is not exist"
	FollowAlreadyExistErrMsg = "follow already exists"
	ContentIsEmptyErrMsg     = "content is empty"

	NotFriendErrMsg  = "not friend"
	NotInGroupErrMsg = "not in group"

	TOTPSecretExpiredErrMsg = "TOTP secret has expired"
	NotGenerateTotpErrMsg   = "TOTP secret not generated or invalidated"
	MFAInvalidCodeErrMsg    = "MFA code is invalid"
)

type ErrNo struct {
	ErrCode int64
	ErrMsg  string
}

func (e ErrNo) Error() string {
	return fmt.Sprintf("err_code=%d, err_msg=%s", e.ErrCode, e.ErrMsg)
}

func NewErrNo(code int64, msg string) ErrNo {
	return ErrNo{code, msg}
}

func (e ErrNo) WithMessage(msg string) ErrNo {
	e.ErrMsg = msg
	return e
}

var (
	Success                         = NewErrNo(SuccessCode, SuccessMsg)
	ServiceErr                      = NewErrNo(ServiceErrCode, ServerErrMsg)
	ParamErr                        = NewErrNo(ParamErrCode, ParamErrMsg)
	UserAlreadyExistErr             = NewErrNo(UserAlreadyExistErrCode, "User already exists")
	AuthorizationFailedErr          = NewErrNo(AuthorizationFailedErrCode, "Authorization failed")
	InvalidFileTypeErr              = NewErrNo(InvalidFileTypeErrCode, InvalidFileTypeErrMsg)
	UserIsNotExistErr               = NewErrNo(UserIsNotExistErrCode, UserIsNotExistErrMsg)
	UserIsNotExistOrPasswordErr     = NewErrNo(AuthorizationFailedErrCode, UserIsNotExistOrPasswordErrMsg)
	FollowRelationAlreadyExistErr   = NewErrNo(FollowRelationAlreadyExistErrCode, "Follow Relation already exist")
	FollowRelationNotExistErr       = NewErrNo(FollowRelationNotExistErrCode, "Follow Relation does not exist")
	FavoriteRelationAlreadyExistErr = NewErrNo(FavoriteRelationAlreadyExistErrCode, "Favorite Relation already exist")
	FavoriteRelationNotExistErr     = NewErrNo(FavoriteRelationNotExistErrCode, "FavoriteRelationNotExistErr")
	FavoriteActionErr               = NewErrNo(FavoriteActionErrCode, FavoriteActionErrMsg)

	MessageAddFailedErr   = NewErrNo(MessageAddFailedErrCode, MessageAddFailedErrMsg)
	FriendIsNotExistErr   = NewErrNo(FriendIsNotExistErrCode, FriendIsNotExistMsg)
	VideoIsNotExistErr    = NewErrNo(VideoIsNotExistErrCode, VideoIsNotExistErrMsg)
	CommentIsNotExistErr  = NewErrNo(CommentIsNotExistErrCode, CommentIsNotExistErrMsg)
	LikeIsNotExistErr     = NewErrNo(LikeIsNotExistErrCode, LikeIsNotExistErrMsg)
	LikeAlreadyExistErr   = NewErrNo(LikeAlreadyExistErrCode, LikeAlreadyExistErrMsg)
	FollowIsNotExistErr   = NewErrNo(FollowIsNotExistErrCode, FollowIsNotExistErrMsg)
	FollowAlreadyExistErr = NewErrNo(FollowAlreadyExistErrCode, FollowAlreadyExistErrMsg)
	ContentIsEmptyErr     = NewErrNo(ContentIsEmptyErrCode, ContentIsEmptyErrMsg)

	NotFriendErr  = NewErrNo(NotFriendErrCode, NotFriendErrMsg)
	NotInGroupErr = NewErrNo(NotInGroupErrCode, NotInGroupErrMsg)

	TOTPSecretExpiredErr = NewErrNo(TOTPSecretExpiredErrCode, TOTPSecretExpiredErrMsg)
	NotGenerateTotpErr   = NewErrNo(NotGenerateTotpErrCode, NotGenerateTotpErrMsg)
	MFAInvalidCodeErr    = NewErrNo(MFAInvalidCodeErrCode, MFAInvalidCodeErrMsg)
)

// ConvertErr convert error to Errno
func ConvertErr(err error) ErrNo {
	Err := ErrNo{}
	if errors.As(err, &Err) {
		return Err
	}

	s := ServiceErr
	s.ErrMsg = err.Error()
	return s
}
