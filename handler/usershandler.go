package handler

import (
	"log/slog"
	"slices"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) UsersInitHandler(c echo.Context) error {
	users := h.GetUsers()
	vctx := view.MakeViewCtx(view.MakeOpts().WithUsers(users).WithOp(constants.OP_CREATE))

	if c.Param("op") == "" {
		return h.renderView(c, vctx.UsersInit())
	} else {
		return h.renderView(c, vctx.UsersLayout())
	}
}

func (h *HandlerContext) Users(c echo.Context, op string) error {
	users := h.GetUsers()
	vctx := view.MakeViewCtx(view.MakeOpts().WithUsers(users).WithOp(op))
	return h.renderView(c, vctx.UsersLayout())
}

func (h *HandlerContext) UserCreateHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")
	hpw, _ := util.HashPassword(password)
	user := entity.User{
		Name:           username,
		HashedPassword: hpw,
		Role:           role,
	}
	err := h.Repo.CreateUser(user)
	if err != nil {
		users := h.GetUsers()
		vctx := view.MakeViewCtx(view.MakeOpts().WithUsers(users).WithSelectedUser(user).WithOp(constants.OP_CREATE).WithErr(err))
		return h.renderView(c, vctx.UsersLayout())
	}

	return h.Users(c, constants.OP_CREATE)
}

func (h *HandlerContext) UserUpdateInitHandler(c echo.Context) error {
	username := c.Param("username")
	user, err := h.Repo.SelectUser(username)
	if err != nil {
		user = entity.User{}
	}
	users := h.GetUsers()
	vctx := view.MakeViewCtx(view.MakeOpts().WithUsers(users).WithSelectedUser(user).WithOp(constants.OP_UPDATE))
	return h.renderView(c, vctx.UsersLayout())
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
		vctx := view.MakeViewCtx(view.MakeOpts().WithSelectedUser(user).WithOp(constants.OP_UPDATE).WithErr(err))
		return h.renderView(c, vctx.UsersLayout())
	}
	return h.Users(c, constants.OP_CREATE)
}

func (h *HandlerContext) UserDeleteHandler(c echo.Context) error {
	username := c.Param("username")
	err := h.Repo.DeleteUser(username)
	if err != nil {
		users := h.GetUsers()
		vctx := view.MakeViewCtx(view.MakeOpts().WithUsers(users).WithOp(constants.OP_CREATE).WithErr(err))
		return h.renderView(c, vctx.UsersLayout())
	}

	return h.Users(c, constants.OP_CREATE)
}

func (h *HandlerContext) ResetPasswordHandler(c echo.Context) error {
	username := c.Param("username")
	newPW := util.GeneratePassword()
	hpw, _ := util.HashPassword(newPW)
	err := h.Repo.ResetPassword(username, hpw)
	if err != nil {
		slog.Error(err.Error(), "error generating new temporary password for ", username)
	}
	users := h.GetUsers()
	vctx := view.MakeViewCtx(view.MakeOpts().WithUsers(users).WithOp(constants.OP_CREATE).WithTempPW(newPW, username))
	return h.renderView(c, vctx.UsersLayout())
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
