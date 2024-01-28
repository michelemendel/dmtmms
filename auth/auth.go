package auth

import (
	"fmt"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	consts "github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/util"
)

type TokenType string

type User struct {
	Name  string
	Token TokenType
}

type Auth struct {
	LoggedInUsers map[TokenType]User
}

func NewUsers() *Auth {
	return &Auth{
		LoggedInUsers: make(map[TokenType]User),
	}
}

func (auth *Auth) Login(c echo.Context, username, password string) error {
	sess, _ := session.Get(consts.SESSION_NAME, c)

	// if (strings.ToLower(username) == "joe" && password == "joe") || (strings.ToLower(username) == "ben" && password == "ben") {
	sess.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 30, // 30 seconds
		// MaxAge:   86400 * 1, // 1 day
		HttpOnly: true,
	}
	newToken := util.GenerateUUID()
	auth.LoggedInUsers[TokenType(newToken)] = User{Name: username, Token: TokenType(newToken)}

	sess.Values[consts.TOKEN_NAME] = newToken
	sess.Save(c.Request(), c.Response())

	fmt.Println("[AUTH]:Login: USER:", username, newToken)
	return nil
	// }

	// fmt.Println("[AUTH]:Login: invalid credentials")
	// return fmt.Errorf("invalid credentials")
}

func (auth *Auth) Logout(c echo.Context) error {
	sess, _ := session.Get(consts.SESSION_NAME, c)
	token := sess.Values[consts.TOKEN_NAME]

	if token != nil {
		user := auth.LoggedInUsers[TokenType(token.(string))]
		fmt.Println("[AUTH]:Logout:", user.Name, token)
	}

	sess.Options.MaxAge = -1
	// sess.
	sess.Save(c.Request(), c.Response())
	delete(auth.LoggedInUsers, TokenType(token.(string)))
	return nil
}

// func (auth *Auth) IsAuthenticated(c echo.Context) bool {
// 	sess, _ := session.Get(consts.SESSION_NAME, c)
// 	token := sess.Values[consts.TOKEN_NAME]
// 	if token != nil {
// 		for _, u := range auth.LoggedInUsers {
// 			if u.Token == token {
// 				fmt.Printf("IsAuthenticated: user:%s, token:%s\n", u.Name, u.Token)
// 				return true
// 			}
// 		}
// 	}
// 	fmt.Println("IsAuthenticated: NOPE")
// 	return false
// }

func (auth *Auth) GetCurrentUser(c echo.Context) (User, error) {
	auth.PrintLoggedInUsers()

	sess, _ := session.Get(consts.SESSION_NAME, c)
	token := sess.Values[consts.TOKEN_NAME]
	if token != nil {
		user := auth.LoggedInUsers[TokenType(token.(string))]
		fmt.Printf("[AUTH]:GetCurrentUser:%s:%s\n", user.Name, user.Token)
		return user, nil
	}
	fmt.Println("[AUTH]:GetCurrentUser: not found")
	return User{}, fmt.Errorf("no user found")
}

func (auth *Auth) GetLoggedInUsers() map[TokenType]User {
	return auth.LoggedInUsers
}

func (auth *Auth) PrintLoggedInUsers() {
	fmt.Println("[AUTH]:--- PrintLoggedInUsers ---")
	for token, user := range auth.LoggedInUsers {
		fmt.Printf("[AUTH]:user:%s:%s\n", user.Name, token)
	}
}
