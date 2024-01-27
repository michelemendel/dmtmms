package handler

import (
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/auth"
	"github.com/michelemendel/dmtmms/constants"
)

func (h *HandlerContext) render(c echo.Context, comp templ.Component) error {
	ctx := context.Background()

	user, err := h.Auth.GetCurrentUser(c)
	fmt.Printf("[RENDER]:user:%s:%s - err:%v\n", user.Name, user.Token, err)
	isLoggedOut := auth.User{} == user
	fmt.Println("[RENDER]:USER:", user, isLoggedOut)

	// if err != nil {
	// 	ctx = context.WithValue(ctx, constants.IsLoggedInKey, false)
	// 	ctx = context.WithValue(ctx, constants.UsernameKey, "error")
	// } else
	if isLoggedOut {
		ctx = context.WithValue(ctx, constants.IsLoggedInKey, false)
		ctx = context.WithValue(ctx, constants.UsernameKey, "anonymous")
	} else {
		ctx = context.WithValue(ctx, constants.IsLoggedInKey, true)
		ctx = context.WithValue(ctx, constants.UsernameKey, user.Name)
	}

	return comp.Render(ctx, c.Response())
}

type CtxKey string

func (h *HandlerContext) SetCtxVal(key CtxKey, value any) {
	h.Ctx = context.WithValue(context.Background(), key, value)
}

func (h *HandlerContext) GetCtxVal(key CtxKey) any {
	return h.Ctx.Value(key)
}
