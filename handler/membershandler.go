package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) MembersHandler(c echo.Context) error {
	fmt.Println("MembersHandler")
	vctx := view.MakeViewCtxDefault()

	members, err := h.Repo.SelectMembers()
	if err != nil {
		return err
	}
	for _, m := range members {
		fmt.Println(m)
	}
	group := entity.Group{}

	return h.renderView(c, vctx.Members(members, group))
}

func (h *HandlerContext) MemberEditHandler(c echo.Context) error {
	return h.renderView(c, view.MemberEdit())
}
