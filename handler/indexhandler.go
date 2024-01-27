package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) IndexHandler(c echo.Context) error {
	fmt.Println("[INDEXHANDLER]")
	return h.render(c, view.Index("THE INDEX"))
}
