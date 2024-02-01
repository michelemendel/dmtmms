package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) PingHandler(c echo.Context) error {
	vctx := view.MakeViewCtx([]entity.User{}, entity.User{}, constants.OP_NONE, nil)
	return h.render(c, vctx.Ping())
}
