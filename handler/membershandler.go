package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) MembersHandler(c echo.Context) error {
	return h.render(c, view.Members(), nil)
}

func (h *HandlerContext) MemberEditHandler(c echo.Context) error {
	return h.render(c, view.MemberEdit(), nil)
}
