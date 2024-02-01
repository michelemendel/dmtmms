package auth

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	consts "github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/util"
)

// TODO: Move to auth package

type Session struct {
	LoggedInUsers map[TokenType]UserSession
}

func NewSession() *Session {
	return &Session{
		LoggedInUsers: make(map[TokenType]UserSession),
	}
}

func (s *Session) Login(c echo.Context, username string) error {
	sess, _ := session.Get(consts.AUTH_SESSION_NAME, c)

	sess.Options = &sessions.Options{
		Path: "/",
		// MaxAge: 30, // 30 seconds
		MaxAge: 3600 * 5, // 5 minutes
		// MaxAge:   86400 * 1, // 1 day
		HttpOnly: true,
	}
	newToken := util.GenerateUUID()
	s.LoggedInUsers[TokenType(newToken)] = UserSession{Name: username, Token: TokenType(newToken)}

	sess.Values[consts.AUTH_TOKEN_NAME] = newToken
	sess.Save(c.Request(), c.Response())

	return nil
}

func (s *Session) Logout(c echo.Context) error {
	sess, _ := session.Get(consts.AUTH_SESSION_NAME, c)
	token := sess.Values[consts.AUTH_TOKEN_NAME]

	if token != nil {
		user := s.LoggedInUsers[TokenType(token.(string))]
		fmt.Println("[SESSION]:Logout:", user.Name, token)
	}

	sess.Options.MaxAge = -1
	sess.Save(c.Request(), c.Response())
	if token != nil {
		delete(s.LoggedInUsers, TokenType(token.(string)))
	}
	return nil
}

func (s *Session) GetCurrentUser(c echo.Context) (UserSession, error) {
	// s.PrintLoggedInUsers()

	sess, _ := session.Get(consts.AUTH_SESSION_NAME, c)
	token := sess.Values[consts.AUTH_TOKEN_NAME]
	if token != nil {
		user := s.LoggedInUsers[TokenType(token.(string))]
		return user, nil
	}
	return UserSession{}, fmt.Errorf("no user found")
}

func (s *Session) GetLoggedInUsers() map[TokenType]UserSession {
	return s.LoggedInUsers
}

func (s *Session) PrintLoggedInUsers() {
	fmt.Println("[SESSION]:--- PrintLoggedInUsers ---")
	for token, user := range s.LoggedInUsers {
		fmt.Printf("[SESSION]:user:%s:%s\n", user.Name, token)
	}
}
