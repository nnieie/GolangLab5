package pack

import (
	"errors"

	"github.com/nnieie/golanglab5/cmd/api/biz/model/base"
	"github.com/nnieie/golanglab5/pkg/errno"
)

func BuildBaseResp(err error) *base.BaseResp {
	if err == nil {
		return &base.BaseResp{
			Code: errno.SuccessCode,
			Msg:  errno.SuccessMsg,
		}
	}
	var errNo errno.ErrNo
	if errors.As(err, &errNo) {
		return &base.BaseResp{
			Code: errNo.ErrCode,
			Msg:  errNo.ErrMsg,
		}
	}
	return &base.BaseResp{
		Code: errno.ServiceErrCode,
		Msg:  err.Error(),
	}
}
