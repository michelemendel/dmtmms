package handler

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

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
	return h.renderView(c, h.ViewCtx.MemberFormModal(entity.Member{}, []string{}, families, groups, constants.OP_CREATE, entity.InputErrors{}))
}

func (h *HandlerContext) MemberCreateHandler(c echo.Context) error {
	inputErrors := entity.NewInputErrors()
	families, _ := h.Repo.SelectFamilies()
	groups, _ := h.Repo.SelectGroups()
	//
	uuid := util.GenerateUUID()
	member, selectedGroupUUIDs := ExtractMemberFromForm(c, uuid)
	inputErrors, areErrors := ValidateInput(member)
	if areErrors {
		return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDs, families, groups, constants.OP_CREATE, inputErrors))
	}
	//
	err := h.Repo.CreateMember(member, selectedGroupUUIDs)
	if err != nil {
		slog.Error(err.Error(), "uuid", uuid, "name", member.Name, "email", member.Email)
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDs, families, groups, constants.OP_CREATE, inputErrors))
	}

	fmt.Println("[CREATE_MEMBER]:", uuid, member.Name, member.Email)

	return h.MembersHandler(c)
}

func ExtractMemberFromForm(c echo.Context, uuid string) (entity.Member, []string) {
	id := c.FormValue("id")
	name := c.FormValue("name")
	dobStr := c.FormValue("dob")
	personnummer := c.FormValue("personnummer")
	email := c.FormValue("email")
	mobile := c.FormValue("mobile")
	groupUUIDsStr := c.FormValue("groups")
	synagogueseat := c.FormValue("synagogueseat")
	membershipFeeTier := c.FormValue("membershipfeetier")
	registeredDateStr := c.FormValue("registereddate")
	deregisteredDateStr := c.FormValue("deregistereddate")
	familyUUID := c.FormValue("family")
	familyGroup := c.FormValue("familygroup")
	//
	dob := util.String2Time(dobStr)
	registeredDate := util.String2Time(registeredDateStr)
	deregisteredDate := util.String2Time(deregisteredDateStr)
	member := entity.NewMember(uuid,
		id, name, dob, personnummer, email,
		mobile, entity.Address{}, synagogueseat, membershipFeeTier,
		registeredDate, deregisteredDate,
		false, false, false, false, entity.MemberStatusActive,
		familyUUID, familyGroup,
	)
	groupUUIDs := strings.Split(groupUUIDsStr, ",")
	return member, groupUUIDs
}
func ValidateInput(member entity.Member) (entity.InputErrors, bool) {
	inputErrors := entity.NewInputErrors()
	//
	if member.Name == "" {
		inputErrors["name"] = entity.NewInputError("name", errors.New("Name is required"))
	}
	if member.ID == "" {
		inputErrors["id"] = entity.NewInputError("name", errors.New("ID is required"))
	}
	areErrors := len(inputErrors) > 0
	return inputErrors, areErrors
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
	inputErrors := entity.NewInputErrors()
	families, _ := h.Repo.SelectFamilies()
	groups, _ := h.Repo.SelectGroups()
	//
	uuid := c.Param("uuid")
	fmt.Println("[UPDATE_MEMBER]: uuid:", uuid)
	member, selectedGroupUUIDs := ExtractMemberFromForm(c, uuid)
	inputErrors, areErrors := ValidateInput(member)
	if areErrors {
		return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDs, families, groups, constants.OP_UPDATE, inputErrors))
	}

	err := h.Repo.UpdateMember(member, selectedGroupUUIDs)
	if err != nil {
		slog.Error(err.Error(), "uuid", uuid, "name", member.Name)
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDs, families, groups, constants.OP_UPDATE, inputErrors))
	}

	return h.MembersHandler(c)
}

func (h *HandlerContext) MemberUpdateHandler(c echo.Context) error {
	return h.MembersHandler(c)
}
