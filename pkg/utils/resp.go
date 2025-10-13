package utils

import (
	"errors"

	"github.com/nnieie/golanglab5/kitex_gen/base"
	"github.com/nnieie/golanglab5/pkg/errno"
)

func BuildBaseResp(err error) *base.BaseResp {
	if err == nil {
		return baseResp(errno.Success)
	}

	e := errno.ErrNo{}
	if errors.As(err, &e) {
		return baseResp(e)
	}

	s := errno.ServiceErr.WithMessage(err.Error())
	return baseResp(s)
}

func baseResp(err errno.ErrNo) *base.BaseResp {
	return &base.BaseResp{
		Code: err.ErrCode,
		Msg:  err.ErrMsg,
	}
}
