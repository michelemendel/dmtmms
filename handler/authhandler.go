package handler

import (
	"fmt"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/util"
)

func (h *HandlerContext) ViewLoginwHandler(c echo.Context) error {
	return h.renderView(c, h.ViewCtx.Login("", nil))
}

func (h *HandlerContext) LoginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.Repo.SelectUser(username)
	if err != nil {
		return h.renderView(c, h.ViewCtx.Login(username, e.ErrInvalidCredentials))
	}

	isAuthed := util.ValidatePassword(password, user.HashedPassword)
	if !isAuthed {
		return h.renderView(c, h.ViewCtx.Login(username, e.ErrInvalidCredentials))
	}

	err = h.Session.Login(c, username)
	if err != nil {
		slog.Error("LoginHandler", "error", e.ErrSystem)
		return h.renderView(c, h.ViewCtx.Login(username, e.ErrSystem))
	}

	c.Response().Header().Set("hx-refresh", "true")
	return h.MembersHandler(c)
}

func (h *HandlerContext) LogoutHandler(c echo.Context) error {
	fmt.Println("[LOGOUTHANDLER]")
	h.Session.Logout(c)
	return h.renderView(c, h.ViewCtx.Login("", nil))
}
