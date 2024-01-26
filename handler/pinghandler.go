package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) PingHandler(c echo.Context) error {
	return h.render(c, view.Ping())
}
