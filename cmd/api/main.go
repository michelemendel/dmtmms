package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/michelemendel/dmtmms/c"
	"github.com/michelemendel/dmtmms/handler"
	"github.com/michelemendel/dmtmms/util"
)

func init() {
	util.InitEnv()
}

func main() {
	strconv.Itoa(4)
	env := os.Getenv(c.APP_ENV)
	logoutput := os.Getenv(c.LOG_OUTPUT)
	webServerPort := os.Getenv(c.WEB_SERVER_PORT)

	fmt.Printf("ENVIRONMENT:\nmode:%s\nlogoutput:%s\nwebServerPort:%s\n", env, logoutput, webServerPort)

	if logoutput == c.LOG_OUTPUT_FILE {
		slog.SetDefault(util.FileLogger())
	} else {
		slog.SetDefault(util.StdOutLogger())
	}

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	// e.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
	// 	fmt.Println("AUTH", username, password)
	// 	if subtle.ConstantTimeCompare([]byte(username), []byte("joe")) == 1 &&
	// 		subtle.ConstantTimeCompare([]byte(password), []byte("secret")) == 1 {
	// 		return true, nil
	// 	}
	// 	return false, nil
	// }))

	e.HTTPErrorHandler = customHTTPErrorHandler

	h := handler.NewHandler()
	Routes(e, h)

	slog.Debug("Starting server", "port", webServerPort)
	e.Logger.Fatal(e.Start(":" + webServerPort))
}

func Routes(e *echo.Echo, h *handler.Handler) {
	e.Static("/public", "public")
	e.GET("/", h.IndexHandler)
	e.GET("/login", h.LoginHandler)
	e.GET("/counts", h.CountsHandler)
	e.POST("/counts", h.CountsHandler)
	e.GET("/ping", h.PingHandler)
}

func customHTTPErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	httpErr, ok := err.(*echo.HTTPError)
	fmt.Println("customHTTPErrorHandler", err, ok, code)
	if ok {
		code = httpErr.Code
	}

	// c.Logger().Error(err)
	// errorPage := fmt.Sprintf("%d.html", code)
	errorPage := "./public/errorPage.html"
	if err := c.File(errorPage); err != nil {
		// c.Logger().Error(err)
		fmt.Println("customHTTPErrorHandler", err)
	}
}
