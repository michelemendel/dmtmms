package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/handler"
)

func Routes(e *echo.Echo, hCtx *handler.HandlerContext) {
	// Static files
	e.Static("/public", "public")

	// Auth
	e.GET(constants.ROUTE_LOGIN, hCtx.ViewLoginwHandler)
	e.POST(constants.ROUTE_LOGIN, hCtx.LoginHandler)
	e.GET(constants.ROUTE_LOGOUT, hCtx.LogoutHandler)

	e.GET(constants.ROUTE_INDEX, hCtx.MembersHandler)

	// Members
	e.GET(constants.ROUTE_MEMBERS, hCtx.MembersHandler)
	e.GET(constants.ROUTE_MEMBER_EDIT, hCtx.MemberEditHandler)

	// Users
	e.GET(constants.ROUTE_USERS, hCtx.UsersInitHandler)
	e.GET(constants.ROUTE_USERS+"/:op", hCtx.UsersInitHandler)
	e.POST(constants.ROUTE_USER_CREATE, hCtx.UserCreateHandler)
	e.GET(constants.ROUTE_USER_UPDATE+"/:username", hCtx.UserUpdateInitHandler)
	e.PUT(constants.ROUTE_USER_UPDATE, hCtx.UserUpdateHandler)
	e.DELETE(constants.ROUTE_USER_DELETE+"/:username", hCtx.UserDeleteHandler)

	//
	e.GET(constants.ROUTE_PING, hCtx.PingHandler)

	// TODO: remove this
	e.GET("/counts", hCtx.CountsHandler)
	e.POST("/counts", hCtx.CountsHandler)
}
