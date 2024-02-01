package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/auth"
	repo "github.com/michelemendel/dmtmms/repository"
)

type HandlerContext struct {
	Ctx     context.Context
	Session *auth.Session
	Repo    *repo.Repo
}

func NewHandlerContext(echo *echo.Echo, auth *auth.Session, repo *repo.Repo) *HandlerContext {
	return &HandlerContext{
		Ctx:     context.Background(),
		Session: auth,
		Repo:    repo,
	}
}
