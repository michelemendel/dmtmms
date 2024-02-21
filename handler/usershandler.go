package handler

import (
	"errors"
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
	return h.renderView(c, h.ViewCtx.Users(users, entity.User{}, "", op, entity.NewInputErrors()))
}

//--------------------------------------------------------------------------------
// Create user

func (h *HandlerContext) UserCreateHandler(c echo.Context) error {
	inputErrors := entity.NewInputErrors()
	users := h.GetUsers()
	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")

	if username == "" || password == "" || role == "" {
		inputErrors["form"] = entity.NewInputError("form", errors.New("username and password are required"))
		return h.renderView(c, h.ViewCtx.Users(users, entity.User{}, "", constants.OP_CREATE, inputErrors))
	}

	hpw, _ := util.HashPassword(password)
	user := entity.User{
		Name:           username,
		HashedPassword: hpw,
		Role:           role,
	}
	err := h.Repo.CreateUser(user)
	if err != nil {
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.Users(users, entity.User{}, "", constants.OP_CREATE, inputErrors))
	}

	return h.Users(c, constants.OP_CREATE)
}

//--------------------------------------------------------------------------------
// Delete user

func (h *HandlerContext) UserDeleteHandler(c echo.Context) error {
	inputErrors := entity.NewInputErrors()
	username := c.Param("username")
	err := h.Repo.DeleteUser(username)
	if err != nil {
		users := h.GetUsers()
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.Users(users, entity.User{}, "", constants.OP_CREATE, inputErrors))
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
	return h.renderView(c, h.ViewCtx.Users(users, user, "", constants.OP_UPDATE, entity.NewInputErrors()))
}

func (h *HandlerContext) UserUpdateHandler(c echo.Context) error {
	inputErrors := entity.NewInputErrors()
	users := h.GetUsers()
	username := c.FormValue("username")
	role := c.FormValue("role")
	user := entity.User{
		Name: username,
		Role: role,
	}

	if username == "" || role == "" {
		inputErrors["form"] = entity.NewInputError("form", errors.New("username and role are required"))
		return h.renderView(c, h.ViewCtx.Users(users, user, "", constants.OP_UPDATE, inputErrors))
	}

	err := h.Repo.UpdateUser(user)
	if err != nil {
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.Users(users, user, "", constants.OP_UPDATE, inputErrors))
	}
	return h.Users(c, constants.OP_CREATE)
}

//--------------------------------------------------------------------------------
// Reset password

func (h *HandlerContext) ResetPasswordHandler(c echo.Context) error {
	inputErrors := entity.NewInputErrors()
	username := c.Param("username")
	tempPW := util.GeneratePassword()
	hpw, _ := util.HashPassword(tempPW)
	err := h.Repo.ResetPassword(username, hpw)
	if err != nil {
		slog.Error(err.Error(), "error generating new temporary password for ", username)
		inputErrors["form"] = entity.NewInputError("form", err)
		return h.renderView(c, h.ViewCtx.Users(h.GetUsers(), entity.User{}, "", constants.OP_CREATE, inputErrors))
	}

	users := h.GetUsers()
	user := entity.User{
		Name: username,
		Role: "",
	}
	return h.renderView(c, h.ViewCtx.Users(users, user, tempPW, constants.OP_CREATE, inputErrors))
}

//--------------------------------------------------------------------------------
// Set password - this has its own GUI

func (h *HandlerContext) SetPasswordInitHandler(c echo.Context) error {
	return h.renderView(c, h.ViewCtx.UserSetPasswordInit())
}

func (h *HandlerContext) SetPasswordHandler(c echo.Context) error {
	inputErrors := entity.NewInputErrors()
	newPW := c.FormValue("newpassword")
	newPWCheck := c.FormValue("newpasswordcheck")
	userSession, err := h.Session.GetLoggedInUser(c)
	username := userSession.Name
	if err != nil {
		slog.Error("error getting current user", "error", err.Error())
		inputErrors["form"] = entity.NewInputError("form", errors.New("there was an error getting the current user"))
		return h.renderView(c, h.ViewCtx.UserSetPassword(newPW, newPWCheck, inputErrors))
	}
	if newPW != newPWCheck {
		inputErrors["form"] = entity.NewInputError("form", errors.New("passwords do not match"))
		return h.renderView(c, h.ViewCtx.UserSetPassword(newPW, newPWCheck, inputErrors))
	}
	hpw, _ := util.HashPassword(newPW)
	err = h.Repo.UpdateUserPassword(entity.User{Name: username, HashedPassword: hpw})
	if err != nil {
		inputErrors["form"] = entity.NewInputError("form", errors.New("there was an error setting the new password"))
		return h.renderView(c, h.ViewCtx.UserSetPassword(newPW, newPWCheck, inputErrors))
	}

	vctx := view.MakeViewCtx(h.Session, view.MakeOpts().WithMsg("password updated successfully"))
	return h.renderView(c, vctx.UserSetPassword("", "", entity.InputErrors{}))
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
