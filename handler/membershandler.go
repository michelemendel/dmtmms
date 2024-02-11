package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/entity"
	repo "github.com/michelemendel/dmtmms/repository"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) MembersHandler(c echo.Context) error {
	members, err := h.MembersFiltered(c)
	if err != nil {
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(err))
		return h.renderView(c, vctx.Members([]entity.Member{}, entity.Member{}, []entity.Group{}, "", "", ""))
	}

	var member entity.Member
	var groups []entity.Group
	// Error here means that the member details are not available, since we haven't selected a member.
	member, groups, err = h.MemberDetails(c)
	if err != nil {
		member = entity.Member{}
		groups = []entity.Group{}
	}

	fromVal, toVal, _ := h.FromTo(c)
	guuid := c.QueryParam("guuid")
	return h.renderView(c, h.ViewCtx.Members(members, member, groups, guuid, fromVal, toVal))
}

func (h *HandlerContext) FromTo(c echo.Context) (string, string, error) {
	fromVal := c.QueryParam("from")
	toVal := c.QueryParam("to")
	return fromVal, toVal, nil
}

func (h *HandlerContext) MembersFiltered(c echo.Context) ([]entity.Member, error) {
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
		return []entity.Member{}, err
	}
	return members, nil
}

func (h *HandlerContext) MemberDetails(c echo.Context) (entity.Member, []entity.Group, error) {
	memberUUID := c.QueryParam("muuid")

	member, err := h.Repo.SelectMemberByUUID(memberUUID)
	if err != nil {
		return entity.Member{}, []entity.Group{}, err
	}

	groups, err := h.Repo.SelectGroupsByMember(memberUUID)
	if err != nil {
		return entity.Member{}, []entity.Group{}, err
	}

	return member, groups, nil
}

//--------------------------------------------------------------------------------
// Create and update member

func (h *HandlerContext) MemberEditHandler(c echo.Context) error {
	return h.renderView(c, h.ViewCtx.MemberEdit())
}
