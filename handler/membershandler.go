package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/entity"
	repo "github.com/michelemendel/dmtmms/repository"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

// queryParams: guuid, from, to
func (h *HandlerContext) MembersHandler(c echo.Context) error {
	return h.Members(c, "Members")
}

func (h *HandlerContext) MembersInnerHandler(c echo.Context) error {
	return h.Members(c, "MembersLayout")
}

func (h *HandlerContext) Members(c echo.Context, component string) error {
	vctx := view.MakeViewCtxDefault()
	guuid := c.QueryParam("guuid")
	fromStr := c.QueryParam("from")
	if fromStr == "" {
		fromStr = "1000-01-01"
	}
	from := util.String2Time(fromStr)
	toStr := c.QueryParam("to")
	if toStr == "" {
		toStr = "3000-01-01"
	}
	to := util.String2Time(toStr)

	filter := repo.MakeFilter(repo.MakeOpts().WithGroupUUID(guuid).WithFrom(from).WithTo(to))

	members, err := h.Repo.SelectMembersByFilter(*filter)
	if err != nil {
		return err
	}
	for _, m := range members {
		fmt.Println(m)
	}
	group := entity.Group{}

	if component == "Members" {
		return h.renderView(c, vctx.Members(members, group))
	} else {
		return h.renderView(c, vctx.MembersLayout(members, group))
	}
}

func (h *HandlerContext) MemberDetailsHandler(c echo.Context) error {
	memberUUID := c.Param("memberuuid")
	vctx := view.MakeViewCtxDefault()

	member, err := h.Repo.SelectMemberByUUID(memberUUID)
	if err != nil {
		vctx := view.MakeViewCtx(view.MakeOpts().WithErr(err))
		return h.renderView(c, vctx.MemberDetails(entity.Member{}, []entity.Group{}))
	}

	groups, err := h.Repo.SelectGroupsByMember(memberUUID)
	if err != nil {
		vctx := view.MakeViewCtx(view.MakeOpts().WithErr(err))
		return h.renderView(c, vctx.MemberDetails(entity.Member{}, []entity.Group{}))
	}

	return h.renderView(c, vctx.MemberDetails(member, groups))
}

func (h *HandlerContext) MemberEditHandler(c echo.Context) error {
	vctx := view.MakeViewCtxDefault()
	return h.renderView(c, vctx.MemberEdit())
}
