package auth

import (
	"github.com/labstack/echo/v4"
)

func (s *Session) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := s.GetLoggedInUser(c)
		path := c.Path()

		if path == "/login" || path == "/node_modules/tw-elements/dist/js*" || path == "/public*" {
			return next(c)
		}

		if sess.Name == "" {
			return echo.ErrUnauthorized
		} else {
			return next(c)
		}
	}
}
