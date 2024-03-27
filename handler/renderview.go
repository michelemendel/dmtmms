package handler

import (
	"context"
	"os"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/auth"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/util"
)

func (h *HandlerContext) renderView(c echo.Context, comp templ.Component) error {
	user, _ := h.Session.GetLoggedInUser(c)
	isLoggedOut := auth.UserSession{} == user

	ctx := context.Background()
	ctx = context.WithValue(ctx, constants.CTX_USER_ROLE_KEY, "")
	if isLoggedOut {
		ctx = context.WithValue(ctx, constants.CTX_USER_NAME_KEY, "anonymous")
	} else {
		u, _ := h.Repo.SelectUser(user.Name)
		ctx = context.WithValue(ctx, constants.CTX_USER_NAME_KEY, user.Name)
		ctx = context.WithValue(ctx, constants.CTX_USER_ROLE_KEY, u.Role)
	}

	// Bypass the login
	if os.Getenv(constants.ENV_BYPASS_LOGIN) == "true" {
		username := "root"
		role := "root"
		h.Session.Login(c, username)
		ctx = context.WithValue(ctx, constants.CTX_USER_NAME_KEY, username)
		ctx = context.WithValue(ctx, constants.CTX_USER_ROLE_KEY, role)
	}

	ctx = context.WithValue(ctx, constants.CTX_IS_XHR_KEY, util.IsHXR(c))

	return comp.Render(ctx, c.Response())
}

type CtxKey string

func (h *HandlerContext) SetCtxVal(key CtxKey, value any) {
	h.Ctx = context.WithValue(context.Background(), key, value)
}

func (h *HandlerContext) GetCtxVal(key CtxKey) any {
	return h.Ctx.Value(key)
}
