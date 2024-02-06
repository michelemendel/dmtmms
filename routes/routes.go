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

	//
	e.GET(constants.ROUTE_INDEX, hCtx.MembersPageHandler)

	// Members
	e.GET(constants.ROUTE_MEMBERS_PAGE, hCtx.MembersPageHandler)
	e.GET(constants.ROUTE_MEMBERS_INTERNAL, hCtx.MembersHandler)
	e.GET(constants.ROUTE_MEMBERS_TABLE, hCtx.MembersTableHandler(false))
	e.GET(constants.ROUTE_MEMBER_EDIT, hCtx.MemberEditHandler)
	e.GET(constants.ROUTE_MEMBER_DETAILS+"/:memberuuid", hCtx.MemberDetailsHandler)

	// Users
	e.GET(constants.ROUTE_USERS, hCtx.UsersHandler(true))
	e.GET(constants.ROUTE_USERS_INTERNAL, hCtx.UsersHandler(false))
	e.GET(constants.ROUTE_USERS+"/:op", hCtx.UsersHandler(false))
	e.POST(constants.ROUTE_USER_CREATE, hCtx.UserCreateHandler)
	e.GET(constants.ROUTE_USER_UPDATE+"/:username", hCtx.UserUpdateInitHandler)
	e.PUT(constants.ROUTE_USER_UPDATE, hCtx.UserUpdateHandler)
	e.DELETE(constants.ROUTE_USER_DELETE+"/:username", hCtx.UserDeleteHandler)
	e.GET(constants.ROUTE_USER_RESET_PW+"/:username", hCtx.ResetPasswordHandler)
	e.GET(constants.ROUTE_USER_SET_PW, hCtx.SetPasswordInitHandler)
	e.POST(constants.ROUTE_USER_SET_PW, hCtx.SetPasswordHandler)

	//
	e.GET(constants.ROUTE_PING, hCtx.PingHandler)
}
