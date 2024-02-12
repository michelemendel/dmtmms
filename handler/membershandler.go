package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/entity"
	repo "github.com/michelemendel/dmtmms/repository"
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

	fromVal, toVal := h.FromTo(c)
	guuid := c.QueryParam("guuid")
	return h.renderView(c, h.ViewCtx.Members(members, member, groups, guuid, fromVal, toVal))
}

func (h *HandlerContext) FromTo(c echo.Context) (string, string) {
	fromVal := c.QueryParam("from")
	toVal := c.QueryParam("to")
	return fromVal, toVal
}

func (h *HandlerContext) MembersFiltered(c echo.Context) ([]entity.Member, error) {
	fmt.Println("MembersFiltered")
	fuuid := c.QueryParam("fuuid")
	guuid := c.QueryParam("guuid")
	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")

	filter := repo.MakeFilter(repo.MakeOpts().WithFamilyUUID(fuuid).WithGroupUUID(guuid).WithFrom(fromStr).WithTo(toStr))
	members, err := h.Repo.SelectMembersByFilter(*filter)
	if err != nil {
		return []entity.Member{}, err
	}
	return members, nil
}

func (h *HandlerContext) MemberDetails(c echo.Context) (entity.Member, []entity.Group, error) {
	fmt.Println("MemberDetails")
	memberUUID := c.QueryParam("muuid")
	if memberUUID == "" {
		return entity.Member{}, nil, nil
	}

	filter := repo.MakeFilter(repo.MakeOpts().WithMemberUUID(memberUUID))
	members, err := h.Repo.SelectMembersByFilter(*filter)
	// member, err := h.Repo.SelectMemberByUUID(memberUUID)
	if err != nil {
		return entity.Member{}, []entity.Group{}, err
	}
	fmt.Println("--- members")
	for _, m := range members {
		fmt.Println(m)
	}

	groups, err := h.Repo.SelectGroupsByMember(memberUUID)
	if err != nil {
		return entity.Member{}, []entity.Group{}, err
	}
	fmt.Println("--- groups")
	for _, g := range groups {
		fmt.Println(g)
	}

	if len(members) > 0 {
		return members[0], groups, nil
	} else {
		return entity.Member{}, groups, nil
	}
}

//--------------------------------------------------------------------------------
// Create and update member

func (h *HandlerContext) MemberEditHandler(c echo.Context) error {
	return h.renderView(c, h.ViewCtx.MemberEdit())
}
