package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/handler"
)

func Routes(e *echo.Echo, hCtx *handler.HandlerContext) {
	// Static files
	e.Static("/public", "public")
	e.Static("/node_modules/tw-elements/dist/js", "node_modules/tw-elements/dist/js/")

	// Auth
	e.POST(constants.ROUTE_LOGIN, hCtx.LoginHandler)
	e.GET(constants.ROUTE_LOGOUT, hCtx.LogoutHandler)

	// Index
	e.GET(constants.ROUTE_INDEX, hCtx.MembersHandler)

	// Users
	e.GET(constants.ROUTE_USERS, hCtx.UsersHandler)
	e.GET(constants.ROUTE_USERS+"/:op", hCtx.UsersHandler)
	e.GET(constants.ROUTE_USER_UPDATE+"/:username", hCtx.UserUpdateInitHandler)
	e.POST(constants.ROUTE_USER_CREATE, hCtx.UserCreateHandler)
	e.PUT(constants.ROUTE_USER_UPDATE, hCtx.UserUpdateHandler)
	e.DELETE(constants.ROUTE_USER_DELETE+"/:username", hCtx.UserDeleteHandler)
	e.GET(constants.ROUTE_USER_RESET_PW+"/:username", hCtx.ResetPasswordHandler)
	e.GET(constants.ROUTE_USER_SET_PW, hCtx.SetPasswordInitHandler)
	e.POST(constants.ROUTE_USER_SET_PW, hCtx.SetPasswordHandler)

	// Members
	e.GET(constants.ROUTE_MEMBERS, hCtx.MembersHandler)
	e.GET(constants.ROUTE_MEMBER_EDIT, hCtx.MemberEditHandler)

	// Groups
	e.GET(constants.ROUTE_GROUPS, hCtx.GroupsHandler)
	e.GET(constants.ROUTE_GROUPS+"/:op", hCtx.GroupsHandler)
	e.POST(constants.ROUTE_GROUP_CREATE, hCtx.GroupCreateHandler)
	e.GET(constants.ROUTE_GROUP_UPDATE+"/:uuid", hCtx.GroupUpdateInitHandler)
	e.PUT(constants.ROUTE_GROUP_UPDATE, hCtx.GroupUpdateHandler)
	e.DELETE(constants.ROUTE_GROUP_DELETE+"/:uuid", hCtx.GroupDeleteHandler)

	// families
	e.GET(constants.ROUTE_FAMILIES, hCtx.FamiliesHandler)
	e.GET(constants.ROUTE_FAMILIES+"/:op", hCtx.FamiliesHandler)
	e.POST(constants.ROUTE_FAMILY_CREATE, hCtx.FamilyCreateHandler)
	e.GET(constants.ROUTE_FAMILY_UPDATE+"/:uuid", hCtx.FamilyUpdateInitHandler)
	e.PUT(constants.ROUTE_FAMILY_UPDATE, hCtx.FamilyUpdateHandler)
	e.DELETE(constants.ROUTE_FAMILY_DELETE+"/:uuid", hCtx.FamilyDeleteHandler)

	//
	e.GET(constants.ROUTE_PING, hCtx.PingHandler)
}
