package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/auth"
	repo "github.com/michelemendel/dmtmms/repository"
	"github.com/michelemendel/dmtmms/view"
)

type HandlerContext struct {
	Ctx     context.Context
	Session *auth.Session
	Repo    *repo.Repo
	ViewCtx *view.ViewCtx
}

func NewHandlerContext(echo *echo.Echo, auth *auth.Session, repo *repo.Repo) *HandlerContext {
	return &HandlerContext{
		Ctx:     context.Background(),
		Session: auth,
		Repo:    repo,
		ViewCtx: view.MakeViewCtxDefault(),
	}
}
