package handler

import (
	"errors"

	// "fmt"
	"log/slog"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/filter"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) MembersHandler(c echo.Context) error {
	nlatest := c.QueryParam("nlatest")
	if nlatest != "" {
		n := util.String2Int(nlatest)
		return h.NLatestMembers(c, n)
	}

	formClose := c.QueryParam("formclose")
	if formClose == "true" {
		// true: We need to refresh the page when the form is closed to make a new search with current filter.
		return h.Members(c, true)
	}
	return h.Members(c, false)
}

func (h *HandlerContext) Members(c echo.Context, doRefresh bool) error {
	members, err := h.MembersFiltered(c)
	if err != nil {
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(err))
		return h.renderView(c, vctx.Members([]entity.Member{}, []entity.Group{}, entity.MemberDetails{}, &filter.Filter{}))
	}

	memberDetails := view.MemberDetailsForPresentation(h.MemberDetails(c))
	if doRefresh {
		c.Response().Header().Set("HX-Refresh", "true")
	}

	return h.renderView(c, h.ViewCtx.Members(members, h.GetGroups(true), memberDetails, h.Filter))
}

func (h *HandlerContext) MembersFiltered(c echo.Context) ([]entity.Member, error) {
	h.Filter.MakeFilterFromQuery(c)
	members, err := h.Repo.SelectMembersByFilter(h.Filter)
	members = h.Filter.SortMembers(c, members)

	if err != nil {
		return []entity.Member{}, err
	}

	return members, nil
}

func (h *HandlerContext) NLatestMembers(c echo.Context, n int) error {
	members, err := h.Repo.SelectNLatestMembers(n)
	if err != nil {
		return h.renderView(c, h.ViewCtx.Members([]entity.Member{}, []entity.Group{}, entity.MemberDetails{}, h.Filter))
	}

	return h.renderView(c, h.ViewCtx.Members(members, []entity.Group{}, entity.MemberDetails{}, &filter.Filter{}))
}

func (h *HandlerContext) MemberDetails(c echo.Context) (entity.Member, []entity.Group) {
	memberUUID := c.QueryParam("muuid")
	if memberUUID == "" {
		return entity.Member{}, []entity.Group{}
	}

	member, err := h.Repo.SelectMember(memberUUID)
	if err != nil {
		return entity.Member{}, []entity.Group{}
	}

	groups, err := h.Repo.SelectGroupsByMember(memberUUID)
	if err != nil {
		return entity.Member{}, []entity.Group{}
	}

	return member, groups
}

//--------------------------------------------------------------------------------
// Create member

func (h *HandlerContext) MemberCreateInitHandler(c echo.Context) error {
	families, _ := h.Repo.SelectFamilies(true)
	groups, _ := h.Repo.SelectGroups(true)
	return h.renderView(c, h.ViewCtx.MemberFormModal(entity.Member{}, []string{}, families, groups, constants.OP_CREATE, h.Filter, entity.InputErrors{}))
}

func (h *HandlerContext) MemberCreateHandler(c echo.Context) error {
	families, _ := h.Repo.SelectFamilies(true)
	groups, _ := h.Repo.SelectGroups(true)

	uuid := util.GenerateUUID()
	member, selectedGroupUUIDs := h.CreatetMemberFromForm(c, uuid)
	inputErrors, areErrors := ValidateInput(member)
	if areErrors {
		c.Response().Header().Set("HX-Retarget", "#memberForm")
		return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDs, families, groups, constants.OP_CREATE, h.Filter, inputErrors))
	}

	err := h.Repo.CreateMember(member, selectedGroupUUIDs)
	if err != nil {
		slog.Error(err.Error(), "uuid", uuid, "name", member.Name, "email", member.Email)
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDs, families, groups, constants.OP_CREATE, h.Filter, inputErrors))
	}

	return h.MembersHandler(c)
}

//--------------------------------------------------------------------------------
// Archive & Delete member

// func (h *HandlerContext) MemberArchiveHandler(c echo.Context) error {
// 	uuid := c.Param("uuid")
// 	err := h.Repo.ArchiveMember(uuid)
// 	if err != nil {
// 		slog.Error(err.Error(), "uuid", uuid)
// 		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErrType(err, view.ErrTypeOnDelete))
// 		return h.renderView(c, vctx.Members([]entity.Member{}, []entity.Group{}, entity.MemberDetails{}, filter.Filter{}))
// 	}
// 	return h.MembersHandler(c)
// }

func (h *HandlerContext) MemberDeleteHandler(c echo.Context) error {
	uuid := c.Param("uuid")
	err := h.Repo.DeleteMember(uuid)
	if err != nil {
		slog.Error(err.Error(), "uuid", uuid)
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErrType(err, view.ErrTypeOnDelete))
		return h.renderView(c, vctx.Members([]entity.Member{}, []entity.Group{}, entity.MemberDetails{}, &filter.Filter{}))
	}
	return h.MembersHandler(c)
}

//--------------------------------------------------------------------------------
// Update member

func (h *HandlerContext) MemberUpdateInitHandler(c echo.Context) error {
	families, _ := h.Repo.SelectFamilies(true)
	groups, _ := h.Repo.SelectGroups(true)
	memberUUID := c.Param("uuid")
	member, _ := h.Repo.SelectMember(memberUUID)
	selectedGroupUUIDs, _ := h.Repo.SelectGroupsByMember(memberUUID)
	selectedGroupUUIDsAsStrings := view.Groups2UUIDsAsStrings(selectedGroupUUIDs)

	// sort := c.QueryParam("sort")
	// order := c.QueryParam("order")
	// c.QueryParams().Set("sort", c.FormValue("sort"))
	// c.QueryParams().Set("order", c.FormValue("order"))
	// c.QueryParams().Set("wsort", "false")
	// c.QueryParams().Set("muuid", memberUUID)
	urlReplacement := h.Filter.URLQuery(memberUUID, "Name", "false")
	// fmt.Println("   [MemberUpdateInitHandler]::", c.Request().URL.Path)
	// fmt.Println("   [MemberUpdateInitHandler]:sort:", sort, order)
	// fmt.Println("   [MemberUpdateInitHandler]:CURR:", urlReplacement)
	c.Response().Header().Set("HX-Replace-Url", urlReplacement)

	return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDsAsStrings, families, groups, constants.OP_UPDATE, h.Filter, entity.InputErrors{}))
}

func (h *HandlerContext) MemberUpdateHandler(c echo.Context) error {
	families, _ := h.Repo.SelectFamilies(true)
	groups, _ := h.Repo.SelectGroups(true)
	memberUUID := c.FormValue("uuid")
	member, selectedGroupUUIDsAsStrings := h.CreatetMemberFromForm(c, memberUUID)

	inputErrors, areErrors := ValidateInput(member)
	if areErrors {
		return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDsAsStrings, families, groups, constants.OP_UPDATE, h.Filter, inputErrors))
	}

	err := h.Repo.UpdateMember(member, selectedGroupUUIDsAsStrings)
	if err != nil {
		slog.Error(err.Error(), "uuid", memberUUID, "name", member.Name)
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.MemberFormModal(member, selectedGroupUUIDsAsStrings, families, groups, constants.OP_UPDATE, h.Filter, inputErrors))
	}

	// To use the current filter, we need to refresh the page to the URL that was set in MemberUpdateInitHandler().
	c.Response().Header().Set("HX-Refresh", "true")
	return h.MembersHandler(c)
}

//--------------------------------------------------------------------------------
// Helper functions

func (h *HandlerContext) CreatetMemberFromForm(c echo.Context, uuid string) (entity.Member, []string) {
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
	receiveHatikvaStr := c.FormValue("receive_hatikvah")
	status := c.FormValue("status")
	familyUUID := c.FormValue("family_uuid")
	familyName := ""
	if familyUUID != "" {
		familyName, _ = h.Repo.GetFamilyNameByUUID(familyUUID)
	}
	//
	dob := util.String2Date(dobStr)
	registeredDate := util.String2Date(registeredDateStr)
	deregisteredDate := util.String2Date(deregisteredDateStr)
	receiveEmail := util.String2Bool(receiveEmailStr)
	receiveMail := util.String2Bool(receiveMailStr)
	receiveHatikva := util.String2Bool(receiveHatikvaStr)
	member := entity.NewMember(uuid,
		util.String2Int(id), name, dob, 0, personnummer, email,
		mobile,
		entity.Address{},
		synagogueseat, membershipFeeTier, registeredDate, deregisteredDate,
		receiveEmail, receiveMail, receiveHatikva, entity.MemberStatus(status),
		familyUUID, familyName,
		time.Time{}, time.Time{},
	)
	params, _ := c.FormParams()
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
