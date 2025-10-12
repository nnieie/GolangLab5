package pack

import (
	"errors"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/nnieie/golanglab5/cmd/api/biz/model/base"
	"github.com/nnieie/golanglab5/pkg/errno"
)

func SendErrResp(c *app.RequestContext, err error) {
	var errNo errno.ErrNo
	ok := errors.As(err, &errNo)
	if !ok {
		c.JSON(consts.StatusOK, base.BaseResp{
			Code: errno.ServiceErrCode,
			Msg:  err.Error(),
		})
		return
	}
	c.JSON(consts.StatusOK, base.BaseResp{
		Code: errNo.ErrCode,
		Msg:  errNo.ErrMsg,
	})
}
