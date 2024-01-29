package handler

import (
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/session"
)

func (h *HandlerContext) render(c echo.Context, comp templ.Component, err error) error {
	ctx := context.Background()

	user, _ := h.Auth.GetCurrentUser(c)
	fmt.Printf("[RENDER]: user:%s, token:%s, err:%v\n", user.Name, user.Token, err)
	isLoggedOut := session.User{} == user
	// fmt.Printf("[RENDER]: isLoggedOut:%v\n", isLoggedOut)

	ctx = context.WithValue(ctx, constants.IS_LOGGEDIN_KEY, false)
	ctx = context.WithValue(ctx, constants.ERROR_KEY, "")
	ctx = context.WithValue(ctx, constants.USER_ROLE_KEY, "")

	if err != nil {
		ctx = context.WithValue(ctx, constants.ERROR_KEY, err.Error())
		ctx = context.WithValue(ctx, constants.USER_NAME_KEY, "error")
	} else if isLoggedOut {
		ctx = context.WithValue(ctx, constants.USER_NAME_KEY, "anonymous")
	} else {
		u, _ := h.Repo.SelectUser(user.Name)
		ctx = context.WithValue(ctx, constants.IS_LOGGEDIN_KEY, true)
		ctx = context.WithValue(ctx, constants.USER_NAME_KEY, user.Name)
		ctx = context.WithValue(ctx, constants.USER_ROLE_KEY, u.Role)
	}

	// TODO: Remove b4to
	// Bypass the auth check for now
	ctx = context.WithValue(ctx, constants.IS_LOGGEDIN_KEY, true)
	ctx = context.WithValue(ctx, constants.USER_NAME_KEY, "root")
	ctx = context.WithValue(ctx, constants.USER_ROLE_KEY, "admin")

	return comp.Render(ctx, c.Response())
}

type CtxKey string

func (h *HandlerContext) SetCtxVal(key CtxKey, value any) {
	h.Ctx = context.WithValue(context.Background(), key, value)
}

func (h *HandlerContext) GetCtxVal(key CtxKey) any {
	return h.Ctx.Value(key)
}
