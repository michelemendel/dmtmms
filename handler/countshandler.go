package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/view"
)

var globalCount int
var sessionCount int

func (h *Handler) CountsHandler(c echo.Context) error {

	val := c.FormValue("val")
	fmt.Println("FormVal:", val)

	if val == "global" {
		globalCount++
	} else if val == "session" {
		sessionCount++
	}

	return render(c, view.Counts(globalCount, sessionCount))
}
