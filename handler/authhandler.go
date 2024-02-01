package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) ViewLoginwHandler(c echo.Context) error {
	vctx := view.MakeViewCtx([]entity.User{}, entity.User{}, constants.OP_NONE, nil)
	return h.render(c, vctx.Login())
}

func (h *HandlerContext) LoginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	// TODO: add logged in user to ctx for the login form.
	vctx := view.MakeViewCtx([]entity.User{}, entity.User{}, constants.OP_NONE, nil)

	user, err := h.Repo.SelectUser(username)
	if err != nil {
		vctx := view.MakeViewCtx([]entity.User{}, entity.User{}, constants.OP_NONE, e.InvalidCredentials)
		return h.render(c, vctx.Login())
	}

	isAuthed := util.ValidatePassword(password, user.HashedPassword)
	fmt.Println("[LOGINHANDLER]:isAuthed:", isAuthed)
	if !isAuthed {
		vctx := view.MakeViewCtx([]entity.User{}, entity.User{}, constants.OP_NONE, e.InvalidCredentials)
		return h.render(c, vctx.Login())
	}

	err = h.Session.Login(c, username)
	if err != nil {
		vctx = view.MakeViewCtx([]entity.User{}, entity.User{}, constants.OP_NONE, e.InvalidCredentials)
	}
	return h.render(c, vctx.Members())
}

func (h *HandlerContext) LogoutHandler(c echo.Context) error {
	fmt.Println("[LOGOUTHANDLER]")
	h.Session.Logout(c)
	vctx := view.MakeViewCtx([]entity.User{}, entity.User{}, constants.OP_NONE, nil)
	return h.render(c, vctx.Login())
}
