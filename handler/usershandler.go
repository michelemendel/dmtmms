package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) UsersHandler(c echo.Context) error {
	users, err := h.Repo.SelectUsers()
	if err != nil {
		users = []entity.User{}
	}
	return h.render(c, view.Users(users), nil)
}

func (h *HandlerContext) ViewUserEditHandler(c echo.Context) error {
	return h.render(c, view.UserEdit(entity.User{}), nil)
}

func (h *HandlerContext) UserEditHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	role := c.FormValue("role")
	fmt.Println("[USEREDITHANDLER]:user:", username, password, role)

	hpw, _ := util.HashPassword(password)
	user := entity.User{
		Name:           username,
		HashedPassword: hpw,
		Role:           role,
	}
	err := h.Repo.CreateUser(user)

	users, err := h.Repo.SelectUsers()
	if err != nil {
		users = []entity.User{}
	}
	return h.render(c, view.Users(users), nil)
}
