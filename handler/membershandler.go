package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) MembersHandler(c echo.Context) error {
	vctx := view.MakeViewCtxDefault()
	return h.renderView(c, vctx.Members())
}

func (h *HandlerContext) MemberEditHandler(c echo.Context) error {
	return h.renderView(c, view.MemberEdit())
}
