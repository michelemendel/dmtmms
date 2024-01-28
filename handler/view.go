package handler

import (
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/auth"
	"github.com/michelemendel/dmtmms/constants"
)

func (h *HandlerContext) render(c echo.Context, comp templ.Component, err error) error {
	ctx := context.Background()

	user, _ := h.Auth.GetCurrentUser(c)
	fmt.Printf("[RENDER]:user:%s:%s - err:%v (%T)\n", user.Name, user.Token, err, err)
	isLoggedOut := auth.User{} == user
	fmt.Println("[RENDER]:isLoggedOut:", isLoggedOut)

	if err != nil {
		ctx = context.WithValue(ctx, constants.ERROR_KEY, err.Error())
		ctx = context.WithValue(ctx, constants.IS_LOGGEDIN_KEY, false)
		ctx = context.WithValue(ctx, constants.USER_NAME_KEY, "error")
	} else if isLoggedOut {
		ctx = context.WithValue(ctx, constants.ERROR_KEY, "")
		ctx = context.WithValue(ctx, constants.IS_LOGGEDIN_KEY, false)
		ctx = context.WithValue(ctx, constants.USER_NAME_KEY, "anonymous")
	} else {
		ctx = context.WithValue(ctx, constants.ERROR_KEY, "")
		ctx = context.WithValue(ctx, constants.IS_LOGGEDIN_KEY, true)
		ctx = context.WithValue(ctx, constants.USER_NAME_KEY, user.Name)
	}

	// TODO: Remove b4to
	// Bypass the auth check for now
	// ctx = context.WithValue(ctx, constants.ERROR_KEY, "")
	// ctx = context.WithValue(ctx, constants.IS_LOGGEDIN_KEY, true)
	// ctx = context.WithValue(ctx, constants.USER_NAME_KEY, "Root")

	return comp.Render(ctx, c.Response())
}

type CtxKey string

func (h *HandlerContext) SetCtxVal(key CtxKey, value any) {
	h.Ctx = context.WithValue(context.Background(), key, value)
}

func (h *HandlerContext) GetCtxVal(key CtxKey) any {
	return h.Ctx.Value(key)
}
