package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

// todo: remove
// func (h *HandlerContext) AuthCheck(c echo.Context) error {
// 	fmt.Println("AuthCheck")
// 	if !h.IsAuthenticated(c) {
// 		fmt.Println("AuthCheck: NOT AUTHENTICATED")
// 		// return c.Redirect(http.StatusMovedPermanently, "/login")
// 		return h.render(c, view.Login())
// 	}
// 	return nil
// }

func (h *HandlerContext) LoginViewHandler(c echo.Context) error {
	return h.render(c, view.Index("THE INDEX"), nil)
}

func (h *HandlerContext) LoginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	user, err := h.Repo.GetUser(username)
	if err != nil {
		return h.render(c, view.Login("Invalid credentials"), fmt.Errorf("invalid credentials"))
	}

	isAuthed := util.ValidatePassword(password, user.HashedPassword)
	fmt.Println("[LOGINHANDLER]:isAuthed:", isAuthed)
	if !isAuthed {
		return h.render(c, view.Login("Invalid credentials"), fmt.Errorf("invalid credentials"))
	}

	err = h.Auth.Login(c, username)
	return h.render(c, view.Index("THE INDEX"), err)
}

func (h *HandlerContext) LogoutHandler(c echo.Context) error {
	fmt.Println("[LOGOUTHANDLER]")
	h.Auth.Logout(c)
	return h.render(c, view.Login(""), nil)
}
