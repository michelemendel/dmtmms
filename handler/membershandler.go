package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/filter"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) MembersHandler(c echo.Context) error {
	f := filter.FilterFromQuery(c)

	members, err := h.MembersFiltered(c, f)
	if err != nil {
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(err))
		return h.renderView(c, vctx.Members([]entity.Member{}, entity.Member{}, []entity.Group{}, filter.Filter{}))
	}

	var member entity.Member
	var groups []entity.Group
	// Error here means that the member details are not available, since we haven't selected a member.
	member, groups, err = h.MemberDetails(c)
	if err != nil {
		member = entity.Member{}
		groups = []entity.Group{}
	}

	return h.renderView(c, h.ViewCtx.Members(members, member, groups, f))
}

func (h *HandlerContext) MembersFiltered(c echo.Context, filter filter.Filter) ([]entity.Member, error) {
	fmt.Println("MembersFiltered: searchTerms:", filter.SearchTerms)

	members, err := h.Repo.SelectMembersByFilter(filter)
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

	filter := filter.MakeFilter(filter.MakeOpts().WithMemberUUID(memberUUID))
	members, err := h.Repo.SelectMembersByFilter(*filter)
	if err != nil {
		return entity.Member{}, []entity.Group{}, err
	}
	for _, m := range members {
		fmt.Println(m)
	}

	groups, err := h.Repo.SelectGroupsByMember(memberUUID)
	if err != nil {
		return entity.Member{}, []entity.Group{}, err
	}
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
