package handler

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/auth"
	"github.com/michelemendel/dmtmms/filter"
	repo "github.com/michelemendel/dmtmms/repository"
	"github.com/michelemendel/dmtmms/view"
)

type HandlerContext struct {
	Ctx     context.Context
	Session *auth.Session
	Repo    *repo.Repo
	ViewCtx *view.ViewCtx
	Filter  *filter.Filter
}

func NewHandlerContext(echo *echo.Echo, session *auth.Session, repo *repo.Repo, f *filter.Filter) *HandlerContext {
	return &HandlerContext{
		Ctx:     context.Background(),
		Session: session,
		Repo:    repo,
		ViewCtx: view.MakeViewCtxDefault(session),
		Filter:  f,
	}
}
