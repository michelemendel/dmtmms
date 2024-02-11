package handler

import (
	"context"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/auth"
	"github.com/michelemendel/dmtmms/constants"
)

func (h *HandlerContext) renderView(c echo.Context, comp templ.Component) error {
	user, _ := h.Session.GetCurrentUser(c)
	isLoggedOut := auth.UserSession{} == user

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CTX_IS_LOGGEDIN_KEY, false)
	ctx = context.WithValue(ctx, constants.CTX_USER_ROLE_KEY, "")
	if isLoggedOut {
		ctx = context.WithValue(ctx, constants.CTX_USER_NAME_KEY, "anonymous")
	} else {
		u, _ := h.Repo.SelectUser(user.Name)
		ctx = context.WithValue(ctx, constants.CTX_IS_LOGGEDIN_KEY, true)
		ctx = context.WithValue(ctx, constants.CTX_USER_NAME_KEY, user.Name)
		ctx = context.WithValue(ctx, constants.CTX_USER_ROLE_KEY, u.Role)
	}

	// TODO: Remove b4to
	// Bypass the auth check for now
	ctx = context.WithValue(ctx, constants.CTX_IS_LOGGEDIN_KEY, true)
	ctx = context.WithValue(ctx, constants.CTX_USER_NAME_KEY, "abe")
	ctx = context.WithValue(ctx, constants.CTX_USER_ROLE_KEY, "admin")

	return comp.Render(ctx, c.Response())
}

type CtxKey string

func (h *HandlerContext) SetCtxVal(key CtxKey, value any) {
	h.Ctx = context.WithValue(context.Background(), key, value)
}

func (h *HandlerContext) GetCtxVal(key CtxKey) any {
	return h.Ctx.Value(key)
}
