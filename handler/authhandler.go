package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
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

func (h *HandlerContext) LoginHandler(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	fmt.Println("[LOGINHANDLER]:", username, password)

	h.Auth.Login(c, username, password)

	return h.render(c, view.Index("THE INDEX"))
}

func (h *HandlerContext) LogoutHandler(c echo.Context) error {
	fmt.Println("[LOGOUTHANDLER]")
	h.Auth.Logout(c)
	return h.render(c, view.Login())
}
