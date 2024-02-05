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
	"github.com/michelemendel/dmtmms/routes"
	"github.com/michelemendel/dmtmms/util"
)

func init() {
	util.InitEnv()
}

func main() {
	env := os.Getenv(consts.ENV_APP_ENV_KEY)
	logoutput := os.Getenv(consts.ENV_LOG_OUTPUT_TYPE_KEY)
	webServerPort := os.Getenv(consts.ENV_WEB_SERVER_PORT_KEY)

	fmt.Printf("ENVIRONMENT:\nmode:%s\nlogoutput:%s\nwebServerPort:%s\n", env, logoutput, webServerPort)

	if logoutput == consts.LOG_OUTPUT_TYPE_FILE {
		slog.SetDefault(util.FileLogger())
	} else {
		slog.SetDefault(util.StdOutLogger())
	}

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(os.Getenv(consts.ENV_SESSION_KEY_KEY)))))
	e.HTTPErrorHandler = customHTTPErrorHandler

	a := auth.NewSession()
	r := repo.NewRepo()
	defer r.DB.Close()

	hCtx := handler.NewHandlerContext(e, a, r)
	routes.Routes(e, hCtx)
	slog.Debug("Starting server", "port", webServerPort)
	e.Logger.Fatal(e.Start(":" + webServerPort))
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	httpErr, ok := err.(*echo.HTTPError)
	if ok {
		code = httpErr.Code
	}

	slog.Warn("httpError", "code", code, "errorMsg", httpErr.Message)

	errorPage := fmt.Sprintf("./public/%d.html", code)
	fileErr := c.File(errorPage)
	if fileErr != nil {
		c.Logger().Error(fileErr)
		fmt.Println("customHTTPErrorHandler", fileErr)
	}
}
