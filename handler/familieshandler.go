package handler

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
)

func (h *HandlerContext) FamiliesHandler(c echo.Context) error {
	return h.Families(c, constants.OP_CREATE)
}

func (h *HandlerContext) Families(c echo.Context, op string) error {
	families := h.GetFamilies(false)
	return h.renderView(c, h.ViewCtx.Families(families, entity.Family{}, op, entity.InputErrors{}))
}

// --------------------------------------------------------------------------------
// Create family

func (h *HandlerContext) FamilyCreateHandler(c echo.Context) error {
	inputErrors := entity.NewInputErrors()
	families := h.GetFamilies(false)
	familyName := c.FormValue("name")

	if familyName == "" {
		inputErrors["form"] = entity.NewInputError("form", errors.New("family name is required"))
		return h.renderView(c, h.ViewCtx.Families(families, entity.Family{}, constants.OP_CREATE, inputErrors))
	}

	family := entity.Family{
		UUID: util.GenerateUUID(),
		Name: familyName,
	}
	err := h.Repo.CreateFamily(family)
	if err != nil {
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.Families(families, family, constants.OP_CREATE, inputErrors))
	}

	return h.Families(c, constants.OP_CREATE)
}

// --------------------------------------------------------------------------------
// Delete family

func (h *HandlerContext) FamilyDeleteHandler(c echo.Context) error {
	inputErrors := entity.NewInputErrors()
	familyUUID := c.Param("uuid")
	err := h.Repo.DeleteFamily(familyUUID)

	if err != nil {
		families := h.GetFamilies(false)
		inputErrors["row"] = entity.NewInputError("row", err)
		return h.renderView(c, h.ViewCtx.Families(families, entity.Family{UUID: familyUUID}, constants.OP_CREATE, inputErrors))
	}

	return h.Families(c, constants.OP_CREATE)
}

// --------------------------------------------------------------------------------
// Update family

func (h *HandlerContext) FamilyUpdateInitHandler(c echo.Context) error {
	families := h.GetFamilies(false)
	familyUUID := c.Param("uuid")
	family, err := h.Repo.SelectFamily(familyUUID)
	if err != nil {
		family = entity.Family{}
	}
	return h.renderView(c, h.ViewCtx.Families(families, family, constants.OP_UPDATE, entity.InputErrors{}))
}

func (h *HandlerContext) FamilyUpdateHandler(c echo.Context) error {
	inputErrors := entity.NewInputErrors()
	families := h.GetFamilies(false)
	familyUUID := c.FormValue("uuid")
	familyName := c.FormValue("name")
	family := entity.Family{
		UUID: familyUUID,
		Name: familyName,
	}

	if familyName == "" {
		inputErrors["form"] = entity.NewInputError("form", errors.New("family name is required"))
		return h.renderView(c, h.ViewCtx.Families(families, family, constants.OP_UPDATE, inputErrors))
	}

	err := h.Repo.UpdateFamily(family)
	if err != nil {
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.Families(families, family, constants.OP_UPDATE, inputErrors))
	}
	return h.Families(c, constants.OP_CREATE)
}

//--------------------------------------------------------------------------------
// Helper functions

func (h *HandlerContext) GetFamilies(withNone bool) []entity.Family {
	families, err := h.Repo.SelectFamilies(withNone)
	if err != nil {
		families = []entity.Family{}
	}
	return families
}
