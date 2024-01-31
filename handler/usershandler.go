package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) UsersInitHandler(c echo.Context) error {
	users := h.GetUsers()
	v := view.MakeUsersViewCtx(users, entity.User{}, constants.OP_CREATE, nil)
	if c.Param("op") == "" {
		return h.render(c, v.UsersInit(), nil)
	} else {
		return h.render(c, v.UsersLayout(), nil)
	}
}

func (h *HandlerContext) Users(c echo.Context, op string) error {
	users := h.GetUsers()
	v := view.MakeUsersViewCtx(users, entity.User{}, op, nil)
	return h.render(c, v.UsersLayout(), nil)
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
		v := view.MakeUsersViewCtx(users, user, constants.OP_CREATE, err)
		return h.render(c, v.UsersLayout(), nil)
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
	v := view.MakeUsersViewCtx(users, user, constants.OP_UPDATE, nil)
	return h.render(c, v.UsersLayout(), nil)
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
		v := view.MakeUsersViewCtx([]entity.User{}, user, constants.OP_UPDATE, err)
		return h.render(c, v.UsersLayout(), err)
	}
	return h.Users(c, constants.OP_CREATE)
}

func (h *HandlerContext) UserDeleteHandler(c echo.Context) error {
	username := c.Param("username")
	err := h.Repo.DeleteUser(username)
	if err != nil {
		users := h.GetUsers()
		v := view.MakeUsersViewCtx(users, entity.User{}, constants.OP_CREATE, err)
		return h.render(c, v.UsersLayout(), err)
	}

	return h.Users(c, constants.OP_CREATE)
}

//--------------------------------------------------------------------------------
// Helper functions

func (h *HandlerContext) GetUsers() []entity.User {
	users, err := h.Repo.SelectUsers()
	if err != nil {
		users = []entity.User{}
	}
	return users
}