package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/entity"
	repo "github.com/michelemendel/dmtmms/repository"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) MembersHandler(c echo.Context) error {
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

	// fmt.Println("fromStr:", fromStr, "toStr:", toStr, "from:", from, "to:", to)

	filter := repo.MakeFilter(repo.MakeOpts().WithGroupUUID(guuid).WithFrom(from).WithTo(to))

	members, err := h.Repo.SelectMembersByFilter(*filter)
	if err != nil {
		return err
	}

	return h.renderView(c, h.ViewCtx.Members(members))
}

func (h *HandlerContext) MemberDetailsHandler(c echo.Context) error {
	memberUUID := c.Param("memberuuid")

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

	return h.renderView(c, h.ViewCtx.MemberDetails(member, groups))
}

//--------------------------------------------------------------------------------
// Create and update member

func (h *HandlerContext) MemberEditHandler(c echo.Context) error {
	return h.renderView(c, h.ViewCtx.MemberEdit())
}
