package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/view"
)

func (h *Handler) LoginHandler(c echo.Context) error {
	return render(c, view.Login())
}
