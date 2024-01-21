package main

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/michelemendel/dmtmms/handler"
)

func main() {
	port := "8080"
	fmt.Println("Listening on :" + port)
	// http.ListenAndServe(":3000", nil)

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	h := handler.NewHandler()
	Routes(e, h)
	e.Logger.Fatal(e.Start(":" + port))
}

// type App struct {
// 	// log *logging.Logger
// }

// func newApp(ctx context.Context, port string) (*App, error) {
// 	app := &App{
// 		// log:       log.Logger("beeguide", logging.RedirectAsJSON(os.Stderr)),
// 	}

// 	return app, nil
// }

func Routes(e *echo.Echo, h *handler.Handler) {
	e.Static("/public", "public")
	e.GET("/ping", h.PingHandler)
	// http.Handle("/", templ.Handler(view.Hello("Milky")))
}
