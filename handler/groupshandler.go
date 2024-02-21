package handler

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
)

func (h *HandlerContext) GroupsHandler(c echo.Context) error {
	return h.Groups(c, constants.OP_CREATE)
}

func (h *HandlerContext) Groups(c echo.Context, op string) error {
	groups := h.GetGroups()
	return h.renderView(c, h.ViewCtx.Groups(groups, entity.Group{}, op, entity.InputErrors{}))
}

//--------------------------------------------------------------------------------
// Create group

func (h *HandlerContext) GroupCreateHandler(c echo.Context) error {
	inputErrors := entity.NewInputErrors()
	groups := h.GetGroups()
	groupName := c.FormValue("name")

	if groupName == "" {
		// vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(fmt.Errorf("group name is required")))
		inputErrors["form"] = entity.NewInputError("form", errors.New("group name is required"))
		return h.renderView(c, h.ViewCtx.Groups(groups, entity.Group{}, constants.OP_CREATE, inputErrors))
	}

	group := entity.Group{
		UUID: util.GenerateUUID(),
		Name: groupName,
	}
	err := h.Repo.CreateGroup(group)
	if err != nil {
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.Groups(groups, group, constants.OP_CREATE, inputErrors))
	}

	return h.Groups(c, constants.OP_CREATE)
}

//--------------------------------------------------------------------------------
// Delete group

func (h *HandlerContext) GroupDeleteHandler(c echo.Context) error {
	inputErrors := entity.NewInputErrors()
	groupUUID := c.Param("uuid")
	err := h.Repo.DeleteGroup(groupUUID)
	if err != nil {
		groups := h.GetGroups()
		inputErrors["row"] = entity.NewInputError("row", err)
		return h.renderView(c, h.ViewCtx.Groups(groups, entity.Group{UUID: groupUUID}, constants.OP_CREATE, inputErrors))
	}

	return h.Groups(c, constants.OP_CREATE)
}

//--------------------------------------------------------------------------------
// Update group

func (h *HandlerContext) GroupUpdateInitHandler(c echo.Context) error {
	groups := h.GetGroups()
	groupUUID := c.Param("uuid")
	group, err := h.Repo.SelectGroup(groupUUID)
	if err != nil {
		group = entity.Group{}
	}
	return h.renderView(c, h.ViewCtx.Groups(groups, group, constants.OP_UPDATE, entity.InputErrors{}))
}

func (h *HandlerContext) GroupUpdateHandler(c echo.Context) error {
	inputErrors := entity.NewInputErrors()
	groups := h.GetGroups()
	groupUUID := c.FormValue("uuid")
	groupName := c.FormValue("name")
	group := entity.Group{
		UUID: groupUUID,
		Name: groupName,
	}

	if groupName == "" {
		inputErrors["form"] = entity.NewInputError("form", errors.New("group name is required"))
		return h.renderView(c, h.ViewCtx.Groups(groups, group, constants.OP_UPDATE, inputErrors))
	}

	err := h.Repo.UpdateGroup(group)
	if err != nil {
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.Groups(groups, group, constants.OP_UPDATE, inputErrors))
	}
	return h.Groups(c, constants.OP_CREATE)
}

//--------------------------------------------------------------------------------
// Helper functions

func (h *HandlerContext) GetGroups() []entity.Group {
	groups, err := h.Repo.SelectGroups()
	if err != nil {
		groups = []entity.Group{}
	}
	return groups
}
