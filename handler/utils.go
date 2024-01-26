package handler

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
)

func (h *HandlerContext) render(c echo.Context, comp templ.Component) error {
	ctx := context.WithValue(context.Background(), constants.UsernameKey, h.GetCtxVal(constants.UsernameKey))
	return comp.Render(ctx, c.Response())
}

type CtxKey string

func (h *HandlerContext) SetCtxVal(key CtxKey, value any) {
	h.Ctx = context.WithValue(context.Background(), key, value)
}

func (h *HandlerContext) GetCtxVal(key CtxKey) any {
	return h.Ctx.Value(key)
}
