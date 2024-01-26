package handler

import (
	"context"

	"github.com/labstack/echo/v4"
)

type HandlerContext struct {
	Ctx           context.Context
	LoggedInUsers []LoggedInUser
}

func NewHandlerContext(echo *echo.Echo) *HandlerContext {
	return &HandlerContext{
		Ctx:           context.Background(),
		LoggedInUsers: make([]LoggedInUser, 0),
	}
}
