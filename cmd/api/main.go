package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/michelemendel/dmtmms/auth"
	consts "github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/handler"
	repo "github.com/michelemendel/dmtmms/repository"
	"github.com/michelemendel/dmtmms/util"
)

func init() {
	util.InitEnv()
}

func main() {
	env := os.Getenv(consts.APP_ENV_KEY)
	logoutput := os.Getenv(consts.LOG_OUTPUT_TYPE_KEY)
	webServerPort := os.Getenv(consts.WEB_SERVER_PORT_KEY)

	fmt.Printf("ENVIRONMENT:\nmode:%s\nlogoutput:%s\nwebServerPort:%s\n", env, logoutput, webServerPort)

	if logoutput == consts.LOG_OUTPUT_TYPE_FILE {
		slog.SetDefault(util.FileLogger())
	} else {
		slog.SetDefault(util.StdOutLogger())
	}

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv(consts.SESSION_KEY_KEY)))))
	e.HTTPErrorHandler = customHTTPErrorHandler

	a := auth.NewUsers()
	r := repo.NewRepo()

	// TODO: Remove b4to
	r.InitDDL()
	r.InitRoot()

	hCtx := handler.NewHandlerContext(e, a, r)
	Routes(e, hCtx)
	slog.Debug("Starting server", "port", webServerPort)
	e.Logger.Fatal(e.Start(":" + webServerPort))
}

func Routes(e *echo.Echo, hCtx *handler.HandlerContext) {
	e.Static("/public", "public")
	e.GET("/login", hCtx.LoginViewHandler)
	e.POST("/login", hCtx.LoginHandler)
	e.GET("/logout", hCtx.LogoutHandler)
	e.GET("/", hCtx.IndexHandler)
	e.GET("/counts", hCtx.CountsHandler)
	e.POST("/counts", hCtx.CountsHandler)
	e.GET("/ping", hCtx.PingHandler)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	httpErr, ok := err.(*echo.HTTPError)
	if ok {
		code = httpErr.Code
	}

	fmt.Printf("customHTTPErrorHandler:ok:%v, code:%v, err:%s\n", ok, code, httpErr.Message)

	// c.Logger().Error(err)
	// errorPage := fmt.Sprintf("%d.html", code)
	// errorPage := "./public/errorPage.html"
	// fileErr := c.File(errorPage)
	// if fileErr != nil {
	// 	c.Logger().Error(fileErr)
	// 	// fmt.Println("customHTTPErrorHandler", fileErr)
	// }
}
