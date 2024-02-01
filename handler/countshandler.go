package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/view"
)

var globalCount int
var sessionCount int

func (h *HandlerContext) CountsHandler(c echo.Context) error {

	val := c.FormValue("val")
	fmt.Println("FormVal:", val)

	if val == "global" {
		globalCount++
	} else if val == "session" {
		sessionCount++
	}

	vctx := view.MakeViewCtxDefault()
	return h.renderView(c, vctx.Counts(globalCount, sessionCount))
}
