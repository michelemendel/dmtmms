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

	// Download
	e.GET("/csv", func(c echo.Context) error {
		msCSV, _ := hCtx.MembersAsCSV(c)
		c.Response().Header().Set("Content-Disposition", "attachment; filename=members.csv")
		return c.Blob(200, "text/csv", msCSV)
	})

	// Auth
	e.POST(constants.ROUTE_LOGIN, hCtx.LoginHandler)
	e.GET(constants.ROUTE_LOGOUT, hCtx.LogoutHandler)

	// Index
	e.GET(constants.ROUTE_INDEX, hCtx.MembersHandler)

	// Users
	e.GET(constants.ROUTE_USERS, hCtx.UsersHandler)
	e.POST(constants.ROUTE_USER_CREATE, hCtx.UserCreateHandler)
	e.DELETE(constants.ROUTE_USER_DELETE+"/:username", hCtx.UserDeleteHandler)
	e.GET(constants.ROUTE_USER_UPDATE+"/:username", hCtx.UserUpdateInitHandler)
	e.PUT(constants.ROUTE_USER_UPDATE, hCtx.UserUpdateHandler)
	// Users, handling passwords
	e.GET(constants.ROUTE_USER_RESET_PW+"/:username", hCtx.ResetPasswordHandler)
	e.GET(constants.ROUTE_USER_SET_PW, hCtx.SetPasswordInitHandler)
	e.POST(constants.ROUTE_USER_SET_PW, hCtx.SetPasswordHandler)

	// Members
	e.GET(constants.ROUTE_MEMBERS, hCtx.MembersHandler)
	e.GET(constants.ROUTE_MEMBER_CREATE, hCtx.MemberCreateInitHandler)
	e.POST(constants.ROUTE_MEMBER_CREATE, hCtx.MemberCreateHandler)
	e.DELETE(constants.ROUTE_MEMBER_DELETE+"/:uuid", hCtx.MemberDeleteHandler)
	e.GET(constants.ROUTE_MEMBER_UPDATE+"/:uuid", hCtx.MemberUpdateInitHandler)
	e.PUT(constants.ROUTE_MEMBER_UPDATE, hCtx.MemberUpdateHandler)

	// Groups
	e.GET(constants.ROUTE_GROUPS, hCtx.GroupsHandler)
	e.POST(constants.ROUTE_GROUP_CREATE, hCtx.GroupCreateHandler)
	e.DELETE(constants.ROUTE_GROUP_DELETE+"/:uuid", hCtx.GroupDeleteHandler)
	e.GET(constants.ROUTE_GROUP_UPDATE+"/:uuid", hCtx.GroupUpdateInitHandler)
	e.PUT(constants.ROUTE_GROUP_UPDATE, hCtx.GroupUpdateHandler)

	// Families
	e.GET(constants.ROUTE_FAMILIES, hCtx.FamiliesHandler)
	e.POST(constants.ROUTE_FAMILY_CREATE, hCtx.FamilyCreateHandler)
	e.DELETE(constants.ROUTE_FAMILY_DELETE+"/:uuid", hCtx.FamilyDeleteHandler)
	e.GET(constants.ROUTE_FAMILY_UPDATE+"/:uuid", hCtx.FamilyUpdateInitHandler)
	e.PUT(constants.ROUTE_FAMILY_UPDATE, hCtx.FamilyUpdateHandler)

	//
	e.GET(constants.ROUTE_PING, hCtx.PingHandler)
}
