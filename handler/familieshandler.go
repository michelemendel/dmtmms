package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) FamiliesHandler(c echo.Context) error {
	return h.Families(c, constants.OP_CREATE)
}

func (h *HandlerContext) Families(c echo.Context, op string) error {
	families := h.GetFamilies()
	return h.renderView(c, h.ViewCtx.Families(families, entity.Family{}, op))
}

func (h *HandlerContext) FamilyCreateHandler(c echo.Context) error {
	families := h.GetFamilies()
	familyGroup := c.FormValue("name")

	if familyGroup == "" {
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(fmt.Errorf("family name is required")))
		return h.renderView(c, vctx.Families(families, entity.Family{}, constants.OP_CREATE))
	}

	family := entity.Family{
		UUID: util.GenerateUUID(),
		Name: familyGroup,
	}
	err := h.Repo.CreateFamily(family)
	fmt.Printf("[CREATE]: %v\n%[1]T\n", err)
	if err != nil {
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErrType(err, view.ErrTypeOnCreate))
		return h.renderView(c, vctx.Families(families, family, constants.OP_CREATE))
	}

	return h.Families(c, constants.OP_CREATE)
}

func (h *HandlerContext) FamilyUpdateInitHandler(c echo.Context) error {
	families := h.GetFamilies()
	familyUUID := c.Param("uuid")
	family, err := h.Repo.SelectFamily(familyUUID)
	if err != nil {
		family = entity.Family{}
	}
	return h.renderView(c, h.ViewCtx.Families(families, family, constants.OP_UPDATE))
}

func (h *HandlerContext) FamilyUpdateHandler(c echo.Context) error {
	familyUUID := c.FormValue("uuid")
	familyGroup := c.FormValue("name")
	family := entity.Family{
		UUID: familyUUID,
		Name: familyGroup,
	}
	err := h.Repo.UpdateFamily(family)
	if err != nil {
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErrType(err, view.ErrTypeOnUpdate))
		return h.renderView(c, vctx.Families([]entity.Family{}, family, constants.OP_UPDATE))
	}
	return h.Families(c, constants.OP_CREATE)
}

func (h *HandlerContext) FamilyDeleteHandler(c echo.Context) error {
	familyUUID := c.Param("uuid")
	err := h.Repo.DeleteFamily(familyUUID)

	fmt.Printf("[DELETE]: %v\n%[1]T\n", err)

	if err != nil {
		families := h.GetFamilies()
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErrType(err, view.ErrTypeOnDelete))
		return h.renderView(c, vctx.Families(families, entity.Family{UUID: familyUUID}, constants.OP_CREATE))
	}

	return h.Families(c, constants.OP_CREATE)
}

//--------------------------------------------------------------------------------
// Helper functions

func (h *HandlerContext) GetFamilies() []entity.Family {
	families, err := h.Repo.SelectFamilies()
	if err != nil {
		families = []entity.Family{}
	}
	return families
}
