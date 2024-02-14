package handler

import (
	"log/slog"

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
		return h.renderView(c, vctx.Members([]entity.Member{}, "", "", "", []entity.MemberDetail{}, []entity.Group{}, filter.Filter{}))
	}

	var member entity.Member
	var groups []entity.Group
	// Error here means that the member details are not available, since we haven't selected a member.
	member, groups, err = h.MemberDetails(c)
	if err != nil {
		member = entity.Member{}
		groups = []entity.Group{}
	}

	memberDatas := entity.GetMemberDetails(member)
	return h.renderView(c, h.ViewCtx.Members(members, member.UUID, member.FamilyUUID, member.FamilyGroup, memberDatas, groups, f))
}

func (h *HandlerContext) MembersFiltered(c echo.Context, filter filter.Filter) ([]entity.Member, error) {
	members, err := h.Repo.SelectMembersByFilter(filter)
	if err != nil {
		return []entity.Member{}, err
	}
	return members, nil
}

func (h *HandlerContext) MemberDetails(c echo.Context) (entity.Member, []entity.Group, error) {
	memberUUID := c.QueryParam("muuid")
	if memberUUID == "" {
		return entity.Member{}, nil, nil
	}

	filter := filter.MakeFilter(filter.MakeOpts().WithMemberUUID(memberUUID))
	members, err := h.Repo.SelectMembersByFilter(*filter)
	if err != nil {
		return entity.Member{}, []entity.Group{}, err
	}

	groups, err := h.Repo.SelectGroupsByMember(memberUUID)
	if err != nil {
		return entity.Member{}, []entity.Group{}, err
	}

	if len(members) > 0 {
		return members[0], groups, nil
	} else {
		return entity.Member{}, groups, nil
	}
}

//--------------------------------------------------------------------------------
// Create member

func (h *HandlerContext) MemberCreateInitHandler(c echo.Context) error {
	return h.renderView(c, h.ViewCtx.MemberForm(""))
}

func (h *HandlerContext) MemberCreateHandler(c echo.Context) error {
	return h.renderView(c, h.ViewCtx.MemberForm(""))
}

//--------------------------------------------------------------------------------
// Delete member

func (h *HandlerContext) MemberDeleteHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	err := h.Repo.DeleteMember(uuid)
	if err != nil {
		slog.Error(err.Error(), "uuid", uuid)
		// vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErrType(err, view.ErrTypeOnDelete))
		// return h.renderView(c, vctx.Members([]entity.Member{}, "", "", "", []entity.MemberDetail{}, []entity.Group{}, filter.Filter{}))
		return h.MembersHandler(c)
	}
	return h.MembersHandler(c)
}

//--------------------------------------------------------------------------------
// Update member

func (h *HandlerContext) MemberUpdateInitHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	return h.renderView(c, h.ViewCtx.MemberForm(uuid))
}

func (h *HandlerContext) MemberUpdateHandler(c echo.Context) error {
	return h.renderView(c, h.ViewCtx.MemberForm(""))
}
