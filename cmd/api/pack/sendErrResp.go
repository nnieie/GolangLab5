package pack

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/nnieie/golanglab5/cmd/api/biz/model/base"
)

func SendErrResp(c *app.RequestContext, err error) {
	c.JSON(consts.StatusOK, struct {
		Base *base.BaseResp `json:"base"`
	}{
		Base: BuildBaseResp(err),
	})
}
