package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	repo "github.com/michelemendel/dmtmms/repository"
	"github.com/michelemendel/dmtmms/session"
)

type HandlerContext struct {
	Ctx  context.Context
	Auth *session.Session
	Repo *repo.Repo
}

func NewHandlerContext(echo *echo.Echo, auth *session.Session, repo *repo.Repo) *HandlerContext {
	return &HandlerContext{
		Ctx:  context.Background(),
		Auth: auth,
		Repo: repo,
	}
}
