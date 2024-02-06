package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) ViewLoginwHandler(c echo.Context) error {
	vctx := view.MakeViewCtxDefault()
	return h.renderView(c, vctx.Login(""))
}

func (h *HandlerContext) LoginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.Repo.SelectUser(username)
	if err != nil {
		vctx := view.MakeViewCtx(view.MakeOpts().WithErr(e.InvalidCredentials))
		return h.renderView(c, vctx.Login(username))
	}

	isAuthed := util.ValidatePassword(password, user.HashedPassword)
	fmt.Println("[LOGINHANDLER]:isAuthed:", isAuthed)
	if !isAuthed {
		vctx := view.MakeViewCtx(view.MakeOpts().WithErr(e.InvalidCredentials))
		return h.renderView(c, vctx.Login(username))
	}

	err = h.Session.Login(c, username)
	if err != nil {
		vctx := view.MakeViewCtx(view.MakeOpts().WithErr(e.InvalidCredentials))
		return h.renderView(c, vctx.AppRoot("", false))
	}

	return h.MembersHandler(c)
}

func (h *HandlerContext) LogoutHandler(c echo.Context) error {
	fmt.Println("[LOGOUTHANDLER]")
	h.Session.Logout(c)
	vctx := view.MakeViewCtxDefault()
	return h.renderView(c, vctx.Login(""))
}
