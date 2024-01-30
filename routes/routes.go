package routes

import (
	"github.com/labstack/echo/v4"
	consts "github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/handler"
)

func Routes(e *echo.Echo, hCtx *handler.HandlerContext) {
	// Static files
	e.Static("/public", "public")

	// Auth
	e.GET(consts.ROUTE_LOGIN, hCtx.ViewLoginwHandler)
	e.POST(consts.ROUTE_LOGIN, hCtx.LoginHandler)
	e.GET(consts.ROUTE_LOGOUT, hCtx.LogoutHandler)

	// Members
	e.GET(consts.ROUTE_INDEX, hCtx.MembersHandler)
	e.GET(consts.ROUTE_MEMBERS, hCtx.MembersHandler)
	e.GET(consts.ROUTE_MEMBER_EDIT, hCtx.MemberEditHandler)

	// Users
	e.GET(consts.ROUTE_USERS, hCtx.UsersNavHandler)
	e.GET(consts.ROUTE_USER_CREATE, hCtx.UserFormHandler)
	e.POST(consts.ROUTE_USER_CREATE, hCtx.UserCreateHandler)
	e.GET(consts.ROUTE_USER_UPDATE+"/:username", hCtx.UserFormHandler)
	e.PUT(consts.ROUTE_USER_UPDATE, hCtx.UserUpdateHandler)
	e.GET(consts.ROUTE_USER_EDIT_CLOSE, hCtx.UsersHandler)
	e.DELETE(consts.ROUTE_USER_DELETE+"/:username", hCtx.UserDeleteHandler)

	//
	e.GET(consts.ROUTE_PING, hCtx.PingHandler)

	// TODO: remove this
	e.GET("/counts", hCtx.CountsHandler)
	e.POST("/counts", hCtx.CountsHandler)
}
