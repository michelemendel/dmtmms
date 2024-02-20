package handler

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/filter"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) MembersHandler(c echo.Context) error {
	f := filter.FilterFromQuery(c)

	members, err := h.MembersFiltered(c, f)
	if err != nil {
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(err))
		return h.renderView(c, vctx.Members([]entity.Member{}, []entity.MemberDetail{}, []entity.Group{}, filter.Filter{}))
	}

	var member entity.Member
	var groups []entity.Group
	// Error here means that the member details are not available, since we haven't selected a member.
	member, groups, err = h.MemberDetails(c)
	if err != nil {
		member = entity.Member{}
		groups = []entity.Group{}
	}

	memberDetails := entity.GetMemberDetails(member)
	return h.renderView(c, h.ViewCtx.Members(members, memberDetails, groups, f))
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
	families, _ := h.Repo.SelectFamilies()
	groups, _ := h.Repo.SelectGroups()
	return h.renderView(c, h.ViewCtx.MemberFormModal(entity.Member{}, families, groups, []string{}, constants.OP_CREATE))
}

func (h *HandlerContext) MemberCreateHandler(c echo.Context) error {
	uuid := util.GenerateUUID()
	id := c.FormValue("id")
	name := c.FormValue("name")
	// dob := c.FormValue("dob")
	// personnummer := c.FormValue("personnummer")
	email := c.FormValue("email")
	// mobile := c.FormValue("mobile")
	// familyUUID := c.FormValue("family")
	// groupUUIDs := c.FormValue("groups")

	fmt.Println("[CREATE_MEMBER]:", uuid, id, name, email)

	return h.MembersHandler(c)
}

//--------------------------------------------------------------------------------
// Archive & Delete member

func (h *HandlerContext) MemberArchiveHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	err := h.Repo.ArchiveMember(uuid)
	if err != nil {
		slog.Error(err.Error(), "uuid", uuid)
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErrType(err, view.ErrTypeOnDelete))
		return h.renderView(c, vctx.Members([]entity.Member{}, []entity.MemberDetail{}, []entity.Group{}, filter.Filter{}))
	}
	return h.MembersHandler(c)
}

func (h *HandlerContext) MemberDeleteHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	err := h.Repo.DeleteMember(uuid)
	if err != nil {
		slog.Error(err.Error(), "uuid", uuid)
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErrType(err, view.ErrTypeOnDelete))
		return h.renderView(c, vctx.Members([]entity.Member{}, []entity.MemberDetail{}, []entity.Group{}, filter.Filter{}))
	}
	return h.MembersHandler(c)
}

//--------------------------------------------------------------------------------
// Update member

func (h *HandlerContext) MemberUpdateInitHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	fmt.Println("[UPDATE_MEMBER]: uuid:", uuid)
	member, err := h.Repo.SelectMemberByUUID(uuid)
	fmt.Println("[UPDATE_MEMBER]:", member.DOB, util.Time2String(member.DOB))
	if err != nil {
		slog.Error(err.Error(), "uuid", uuid)
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErrType(err, view.ErrTypeOnUpdate))
		return h.renderView(c, vctx.Members([]entity.Member{}, []entity.MemberDetail{}, []entity.Group{}, filter.Filter{}))
	}
	families, _ := h.Repo.SelectFamilies()
	groups, _ := h.Repo.SelectGroups()
	selectedGroupUUIDs, _ := h.Repo.SelectGroupUUIDsByMember(member.UUID)
	return h.renderView(c, h.ViewCtx.MemberFormModal(member, families, groups, selectedGroupUUIDs, constants.OP_UPDATE))
}

func (h *HandlerContext) MemberUpdateHandler(c echo.Context) error {
	return h.MembersHandler(c)
}
