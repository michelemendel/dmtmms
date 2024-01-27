package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/auth"
)

type HandlerContext struct {
	Ctx  context.Context
	Auth *auth.Auth
}

func NewHandlerContext(echo *echo.Echo, auth *auth.Auth) *HandlerContext {
	return &HandlerContext{
		Ctx:  context.Background(),
		Auth: auth,
	}
}
