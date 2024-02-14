package handler

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) UsersHandler(c echo.Context) error {
	return h.Users(c, constants.OP_CREATE)
}

func (h *HandlerContext) Users(c echo.Context, op string) error {
	users := h.GetUsers()
	return h.renderView(c, h.ViewCtx.Users(users, entity.User{}, op))
}

//--------------------------------------------------------------------------------
// Create user

func (h *HandlerContext) UserCreateHandler(c echo.Context) error {
	users := h.GetUsers()
	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")

	if username == "" || password == "" || role == "" {
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(fmt.Errorf("username, password, and role are required")))
		return h.renderView(c, vctx.Users(users, entity.User{}, constants.OP_CREATE))
	}

	hpw, _ := util.HashPassword(password)
	user := entity.User{
		Name:           username,
		HashedPassword: hpw,
		Role:           role,
	}
	err := h.Repo.CreateUser(user)
	if err != nil {
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(err))
		return h.renderView(c, vctx.Users(users, user, constants.OP_CREATE))
	}

	return h.Users(c, constants.OP_CREATE)
}

//--------------------------------------------------------------------------------
// Delete user

func (h *HandlerContext) UserDeleteHandler(c echo.Context) error {
	username := c.Param("username")
	err := h.Repo.DeleteUser(username)
	if err != nil {
		users := h.GetUsers()
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(err))
		return h.renderView(c, vctx.Users(users, entity.User{}, constants.OP_CREATE))
	}

	return h.Users(c, constants.OP_CREATE)
}

//--------------------------------------------------------------------------------
// Update user

func (h *HandlerContext) UserUpdateInitHandler(c echo.Context) error {
	users := h.GetUsers()
	username := c.Param("username")
	user, err := h.Repo.SelectUser(username)
	if err != nil {
		user = entity.User{}
	}
	return h.renderView(c, h.ViewCtx.Users(users, user, constants.OP_UPDATE))
}

func (h *HandlerContext) UserUpdateHandler(c echo.Context) error {
	username := c.FormValue("username")
	role := c.FormValue("role")
	user := entity.User{
		Name: username,
		Role: role,
	}
	err := h.Repo.UpdateUser(user)
	if err != nil {
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(err))
		return h.renderView(c, vctx.Users([]entity.User{}, user, constants.OP_UPDATE))
	}
	return h.Users(c, constants.OP_CREATE)
}

//--------------------------------------------------------------------------------
// Reset and set password

func (h *HandlerContext) ResetPasswordHandler(c echo.Context) error {
	username := c.Param("username")
	newPW := util.GeneratePassword()
	hpw, _ := util.HashPassword(newPW)
	err := h.Repo.ResetPassword(username, hpw)
	if err != nil {
		slog.Error(err.Error(), "error generating new temporary password for ", username)
	}

	users := h.GetUsers()
	vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithTempPW(newPW, username))
	return h.renderView(c, vctx.Users(users, entity.User{}, constants.OP_CREATE))
}

func (h *HandlerContext) SetPasswordInitHandler(c echo.Context) error {
	return h.renderView(c, h.ViewCtx.UserSetPasswordInit())
}

func (h *HandlerContext) SetPasswordHandler(c echo.Context) error {
	newPW := c.FormValue("newpassword")
	newPWCheck := c.FormValue("newpasswordcheck")
	userSession, err := h.Session.GetCurrentUser(c)
	username := userSession.Name
	if err != nil {
		slog.Error("error getting current user", "error", err.Error())
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(fmt.Errorf("there was an error getting the current user")))
		return h.renderView(c, vctx.UserSetPassword(newPW, newPWCheck))
	}
	if newPW != newPWCheck {
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(fmt.Errorf("passwords do not match")))
		return h.renderView(c, vctx.UserSetPassword(newPW, newPWCheck))
	}
	hpw, _ := util.HashPassword(newPW)
	err = h.Repo.UpdateUserPassword(entity.User{Name: username, HashedPassword: hpw})
	if err != nil {
		vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithErr(fmt.Errorf("there was an error setting the new password")))
		return h.renderView(c, vctx.UserSetPassword(newPW, newPWCheck))
	}

	vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithMsg("password updated successfully"))
	return h.renderView(c, vctx.UserSetPassword("", ""))
}

//--------------------------------------------------------------------------------
// Helper functions

func (h *HandlerContext) GetUsers() []entity.User {
	users, err := h.Repo.SelectUsers()
	if err != nil {
		users = []entity.User{}
	}

	// We don't want to allow the root user to be edited or deleted,
	// so we remove him from the list.
	for i, user := range users {
		if user.Name == "root" {
			users = slices.Delete(users, i, 1)
		}
	}

	return users
}
