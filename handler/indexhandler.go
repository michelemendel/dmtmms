package handler

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/view"
)

func (h *HandlerContext) IndexHandler(c echo.Context) error {
	fmt.Println("IndexHandler")
	h.AuthCheck(c)
	return render(c, view.Index("THE INDEX"))
}
