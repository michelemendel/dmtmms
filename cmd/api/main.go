package main

import (
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/michelemendel/dmtmms/handler"
	"github.com/michelemendel/dmtmms/util"
)

func init() {
	util.InitEnv()
}

func main() {
	// slog.SetDefault(util.StdOutLogger())
	slog.SetDefault(util.FileLogger())

	port := os.Getenv("PORT")

	slog.Debug("A:---")
	slog.Info("B:---")

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	h := handler.NewHandler()
	Routes(e, h)

	slog.Debug("Starting server", "port", port)

	e.Logger.Fatal(e.Start(":" + port))
}

func Routes(e *echo.Echo, h *handler.Handler) {
	e.Static("/public", "public")
	e.GET("/ping", h.PingHandler)
	// http.Handle("/", templ.Handler(view.Hello("Milky")))
}
