package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) PingHandler(c echo.Context) error {
	vctx := view.MakeViewCtxDefault()
	return h.renderView(c, vctx.Ping())
}
