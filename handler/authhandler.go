package handler

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	consts "github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/util"
	"github.com/michelemendel/dmtmms/view"
)

type LoggedInUser struct {
	name  string
	token string
}

func (h *HandlerContext) IsAuthenticated(c echo.Context) bool {
	sess, _ := session.Get(consts.SESSION_NAME, c)
	token := sess.Values[consts.TOKEN_NAME]
	if token != nil {
		for _, u := range h.loggedInUsers {
			if u.token == token {
				fmt.Printf("IsAuthenticated: user:%s, token:%s\n", u.name, u.token)
				return true
			}
		}
	}
	fmt.Println("IsAuthenticated: NOPE")
	return false
}

func (h *HandlerContext) AuthCheck(c echo.Context) error {
	fmt.Println("AuthCheck")
	if !h.IsAuthenticated(c) {
		fmt.Println("AuthCheck: NOT AUTHENTICATED")
		// return c.Redirect(http.StatusMovedPermanently, "/login")
		return render(c, view.Login())
	}
	return nil
}

// func (h *HandlerContext) Auth(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		fmt.Println("[auth_middleware] START: currRoute:", c.Get("currentRoute"))
// 		if !h.IsAuthenticated(c) {
// 			fmt.Println("[auth_middleware] NOT AUTHENTICATED")
// 			// c.Redirect(http.StatusMovedPermanently, "/login")
// 			// render(c, view.Login())
// 			// return nil
// 		}

// 		fmt.Println("[auth_middleware] NEXT")
// 		if err := next(c); err != nil {
// 			c.Error(err)
// 		}
// 		fmt.Println("[auth_middleware] AFTER NEXT")
// 		return nil
// 	}
// }

func (h *HandlerContext) LoginHandler(c echo.Context) error {
	if h.IsAuthenticated(c) {
		// return c.Redirect(http.StatusMovedPermanently, "/")
		return render(c, view.Index("THE INDEX"))
	}

	sess, _ := session.Get(consts.SESSION_NAME, c)
	username := c.FormValue("username")
	password := c.FormValue("password")

	fmt.Println("LoginHandler: un/pw:", username, password)

	if username == "joe" && password == "joe" {
		sess.Options = &sessions.Options{
			Path:   "/",
			MaxAge: 30, // 30 seconds
			// MaxAge:   86400 * 7, // 7 days
			HttpOnly: true,
		}
		newToken := util.GenerateUuid()
		h.loggedInUsers = append(h.loggedInUsers, LoggedInUser{name: username, token: newToken})
		sess.Values[consts.TOKEN_NAME] = newToken
		sess.Save(c.Request(), c.Response())
		// return c.Redirect(http.StatusMovedPermanently, "/")
		fmt.Println("LoginHandler: A")
		return render(c, view.Index("THE INDEX"))
	}

	fmt.Println("LoginHandler: B")
	return render(c, view.Login())
}

func (h *HandlerContext) LogoutHandler(c echo.Context) error {
	fmt.Println("LogoutHandler")
	sess, _ := session.Get(consts.SESSION_NAME, c)
	fmt.Println("LogoutHandler:A:", sess.Values)
	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())
	fmt.Println("LogoutHandler:B:", sess.Values)
	// return c.Redirect(http.StatusMovedPermanently, "/login")
	return render(c, view.Login())
}
