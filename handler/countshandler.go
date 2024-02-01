package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
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

	vctx := view.MakeViewCtx([]entity.User{}, entity.User{}, constants.OP_NONE, nil)
	return h.render(c, vctx.Counts(globalCount, sessionCount))
}
