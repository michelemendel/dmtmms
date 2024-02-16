package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/util"
)

func (h *HandlerContext) ViewLoginHandler(c echo.Context) error {
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

	h.Session.Login(c, username)

	return c.Redirect(302, "/members?l=ok")
	// return h.renderView(c, h.ViewCtx.Members([]entity.Member{}, "", "", "", []entity.MemberDetail{}, []entity.Group{}, filter.Filter{}))
}

func (h *HandlerContext) LogoutHandler(c echo.Context) error {
	h.Session.Logout(c)
	return h.renderView(c, h.ViewCtx.Login("", nil))
}
