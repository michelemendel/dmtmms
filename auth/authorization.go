package auth

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
)

func (s *Session) IsAuthorized(userRole string, path string) bool {

	if path == constants.ROUTE_INDEX ||
		path == constants.ROUTE_LOGIN ||
		path == constants.ROUTE_LOGOUT ||
		path == constants.ROUTE_USER_SET_PW ||
		path == constants.ROUTE_MEMBERS {
		return true
	}

	switch userRole {
	case "root":
		return true
	case "admin":
		switch path {
		case constants.ROUTE_FAMILIES:
			return true
		case constants.ROUTE_GROUPS:
			return true
		case constants.ROUTE_USERS:
			return true
		}
	case "edit":
		switch path {
		case constants.ROUTE_FAMILIES:
			return true
		case constants.ROUTE_GROUPS:
			return true
		case constants.ROUTE_USERS:
			return false
		}
	case "read":
		switch path {
		case constants.ROUTE_FAMILIES:
			return false
		case constants.ROUTE_GROUPS:
			return false
		case constants.ROUTE_USERS:
			return false
		}
	}
	return false
}

// Authorize is a middleware to check if the user is authorized to access the route
func (s *Session) Authorize(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := s.GetCurrentUser(c)
		path := c.Path()

		if sess.Name == "" || path == "/node_modules/tw-elements/dist/js*" || path == "/public*" {
			return next(c)
		}

		user, _ := s.Repo.SelectUser(sess.Name)

		fmt.Printf("Authorize: role:%s, path:%s\n", user.Role, path)
		if s.IsAuthorized(user.Role, path) {
			return next(c)
		} else {
			return next(c)
			// return echo.ErrUnauthorized
		}
	}
}
