package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) GetUsers() []entity.User {
	users, err := h.Repo.SelectUsers()
	if err != nil {
		users = []entity.User{}
	}
	return users
}

func (h *HandlerContext) UsersNavHandler(c echo.Context) error {
	users := h.GetUsers()
	return h.render(c, view.UsersNav(users), nil)
}

func (h *HandlerContext) UsersHandler(c echo.Context) error {
	users := h.GetUsers()
	return h.render(c, view.Users(users), nil)
}

func (h *HandlerContext) UserFormHandler(c echo.Context) error {
	username := c.Param("username")
	fmt.Println("[USEREDITFORM]:username:", username)
	user, err := h.Repo.SelectUser(username)
	fmt.Println("[USEREDITFORM]:user:", user)
	if err != nil {
		user = entity.User{}
	}
	op := constants.OP_CREATE
	if user.Name != "" {
		op = constants.OP_UPDATE
	}
	return h.render(c, view.UserFormType(user, op), nil)
}

func (h *HandlerContext) UserCreateHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")
	fmt.Println("[USER_EDIT]:user:", username, password, role)

	hpw, _ := util.HashPassword(password)
	user := entity.User{
		Name:           username,
		HashedPassword: hpw,
		Role:           role,
	}
	err := h.Repo.CreateUser(user)
	if err != nil {
		fmt.Println("[USER_CREATE]:error:", err)
		users := h.GetUsers()
		return h.render(c, view.Users(users), err)
		// return h.render(c, view.UserEdit(user, constants.OP_CREATE), err)
		// return h.render(c, view.UserForm(user), err)
		// return h.render(c, view.UserFormError(err), nil)
	}

	return h.UsersHandler(c)
}

func (h *HandlerContext) UserUpdateHandler(c echo.Context) error {
	username := c.FormValue("username")
	role := c.FormValue("role")
	fmt.Println("[USER_UPDATE]:user:", username, role)

	user := entity.User{
		Name: username,
		Role: role,
	}

	err := h.Repo.UpdateUser(user)
	if err != nil {
		return h.render(c, view.UserFormType(user, constants.OP_UPDATE), err)
	}

	return h.UsersHandler(c)
}

func (h *HandlerContext) UserDeleteHandler(c echo.Context) error {
	username := c.Param("username")
	fmt.Println("[USER_DELETE]:user:", username)

	err := h.Repo.DeleteUser(username)
	if err != nil {
		// return h.UsersHandler(c)
		users := h.GetUsers()
		return h.render(c, view.Users(users), err)

	}

	return h.UsersHandler(c)
}
