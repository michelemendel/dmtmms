package handler

import (
	"errors"
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

	memberDetails := entity.GetMemberDetailsForPresentation(member)
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
	families, _ := h.Repo.SelectFamilies()
	groups, _ := h.Repo.SelectGroups()
	//
	uuid := util.GenerateUUID()
	member, selectedGroupUUIDs := h.CreatetMemberFromForm(c, uuid)
	inputErrors, areErrors := ValidateInput(member)
	if areErrors {
		c.Response().Header().Set("HX-Retarget", "#memberForm")
		return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDs, families, groups, constants.OP_CREATE, inputErrors))
	}
	//
	err := h.Repo.CreateMember(member, selectedGroupUUIDs)
	if err != nil {
		slog.Error(err.Error(), "uuid", uuid, "name", member.Name, "email", member.Email)
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDs, families, groups, constants.OP_CREATE, inputErrors))
	}

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
	families, _ := h.Repo.SelectFamilies()
	groups, _ := h.Repo.SelectGroups()
	memberUUID := c.Param("uuid")
	filter := filter.MakeFilter(filter.MakeOpts().WithMemberUUID(memberUUID))
	member := entity.Member{}
	members, _ := h.Repo.SelectMembersByFilter(*filter)
	selectedGroupUUIDsAsStrings := []string{}
	selectedGroupUUIDs, _ := h.Repo.SelectGroupsByMember(memberUUID)
	if len(members) > 0 {
		member = members[0]
		selectedGroupUUIDsAsStrings = entity.Groups2UUIDsAsStrings(selectedGroupUUIDs)
	}

	return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDsAsStrings, families, groups, constants.OP_UPDATE, entity.InputErrors{}))
}

func (h *HandlerContext) MemberUpdateHandler(c echo.Context) error {
	families, _ := h.Repo.SelectFamilies()
	groups, _ := h.Repo.SelectGroups()
	uuid := c.FormValue("uuid")
	member, selectedGroupUUIDsAsStrings := h.CreatetMemberFromForm(c, uuid)

	inputErrors, areErrors := ValidateInput(member)
	if areErrors {
		return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDsAsStrings, families, groups, constants.OP_UPDATE, inputErrors))
	}

	err := h.Repo.UpdateMember(member, selectedGroupUUIDsAsStrings)
	if err != nil {
		slog.Error(err.Error(), "uuid", uuid, "name", member.Name)
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDsAsStrings, families, groups, constants.OP_UPDATE, inputErrors))
	}
	return h.MembersHandler(c)
}

//--------------------------------------------------------------------------------
// Helper functions

func (h *HandlerContext) CreatetMemberFromForm(c echo.Context, uuid string) (entity.Member, []string) {
	params, _ := c.FormParams()
	id := c.FormValue("id")
	name := c.FormValue("name")
	dobStr := c.FormValue("dob")
	personnummer := c.FormValue("personnummer")
	email := c.FormValue("email")
	mobile := c.FormValue("mobile")
	synagogueseat := c.FormValue("synagogue_seat")
	membershipFeeTier := c.FormValue("membership_fee_tier")
	registeredDateStr := c.FormValue("registered_date")
	deregisteredDateStr := c.FormValue("deregistered_date")
	receiveEmailStr := c.FormValue("receive_email")
	receiveMailStr := c.FormValue("receive_mail")
	receiveHatikvaStr := c.FormValue("receive_hatikva")
	status := c.FormValue("status")
	familyUUID := c.FormValue("family_uuid")
	familyName := ""
	if familyUUID != "" {
		familyName, _ = h.Repo.GetFamilyNameByUUID(familyUUID)
	}
	//
	dob := util.String2Time(dobStr)
	registeredDate := util.String2Time(registeredDateStr)
	deregisteredDate := util.String2Time(deregisteredDateStr)
	receiveEmail := util.String2Bool(receiveEmailStr)
	receiveMail := util.String2Bool(receiveMailStr)
	receiveHatikva := util.String2Bool(receiveHatikvaStr)
	member := entity.NewMember(uuid,
		id, name, dob, personnummer, email,
		mobile,
		entity.Address{},
		synagogueseat, membershipFeeTier, registeredDate, deregisteredDate,
		receiveEmail, receiveMail, receiveHatikva, false, entity.MemberStatus(status),
		familyUUID, familyName,
	)
	groupUUIDs := params["groups"]
	return member, groupUUIDs
}
func ValidateInput(member entity.Member) (entity.InputErrors, bool) {
	inputErrors := entity.NewInputErrors()
	//
	if member.Name == "" {
		inputErrors["name"] = entity.NewInputError("name", errors.New("name is required"))
	}
	areErrors := len(inputErrors) > 0
	return inputErrors, areErrors
}
