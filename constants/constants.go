package constants

const (
	DATE_TIME_FRMT = "2006-01-02 15:04:05"
	DATE_FRMT      = "2006-01-02"
	DATE_MIN       = "0001-01-01"
	DATE_MAX       = "3000-01-01"
)

// Env keys
const (
	ENV_APP_ENV_KEY         = "APP_ENV"
	ENV_FILE_NAME_KEY       = "LOG_FILE_NAME"
	ENV_LOG_OUTPUT_TYPE_KEY = "LOG_OUTPUT_TYPE"
	ENV_WEB_SERVER_PORT_KEY = "WEB_SERVER_PORT"
	ENV_SESSION_KEY_KEY     = "SESSION_KEY"
	ENV_DEV_DB_DIR_KEY      = "DEV_DB_DIR"
	ENV_PROD_DB_DIR_KEY     = "PROD_DB_DIR"
	ENV_DB_NAME_KEY         = "DB_NAME"
	ENV_BYPASS_LOGIN        = "BYPASS_LOGIN"
)

// Log output types
const (
	LOG_OUTPUT_TYPE_FILE   = "file"
	LOG_OUTPUT_TYPE_STDOUT = "stdout"
)

// Auth keys
const (
	AUTH_SESSION_NAME = "session"
	AUTH_TOKEN_NAME   = "token"
	// auth objects
	AUTH_NAV_USERS = "navUsers"
)

// Context keys
const (
	CTX_USER_NAME_KEY   = "username"
	CTX_USER_ROLE_KEY   = "userrole"
	CTX_IS_LOGGEDIN_KEY = "isloggedinkey"
	CTX_IS_XHR_KEY      = "isxhr"
)

// Routes
const (
	ROUTE_ANY = "any"
	//
	ROUTE_LOGIN  = "/login"
	ROUTE_LOGOUT = "/logout"
	//
	ROUTE_INDEX = "/"
	//
	// Users
	ROUTE_USERS         = "/users"
	ROUTE_USER_CREATE   = "/user/create"
	ROUTE_USER_DELETE   = "/user/delete"
	ROUTE_USER_UPDATE   = "/user/update"
	ROUTE_USER_RESET_PW = "/user/resetpw"
	ROUTE_USER_SET_PW   = "/user/setpw"

	//
	// Members
	ROUTE_MEMBERS        = "/members"
	ROUTE_MEMBER_DETAILS = "/memberdetails"
	ROUTE_MEMBER_CREATE  = "/member/create"
	ROUTE_MEMBER_DELETE  = "/member/delete"
	ROUTE_MEMBER_UPDATE  = "/member/update"

	//
	// Groups
	ROUTE_GROUPS       = "/groups"
	ROUTE_GROUP_CREATE = "/group/create"
	ROUTE_GROUP_DELETE = "/group/delete"
	ROUTE_GROUP_UPDATE = "/group/update"

	//
	// Families
	ROUTE_FAMILIES      = "/families"
	ROUTE_FAMILY_CREATE = "/family/create"
	ROUTE_FAMILY_DELETE = "/family/delete"
	ROUTE_FAMILY_UPDATE = "/family/update"

	//
	ROUTE_DOWNLOAD = "/download"

	//
	ROUTE_PING = "/ping"
)

// View operations
const (
	OP_NONE   = ""
	OP_CREATE = "create"
	OP_UPDATE = "update"
)

// Roles
const (
	ROLE_ROOT  = "root"
	ROLE_ADMIN = "admin"
	ROLE_EDIT  = "edit"
	ROLE_READ  = "read"
)

// CSS
const (
	ClassIndeterminate = "relative peer shrink-0 appearance-none w-4 h-4 text-red-600 border-2 border-blue-500 rounded-sm bg-white mt-1 checked:bg-red-600 checked:border-0 focus:outline-none focus:ring-offset-0 focus:ring-2 focus:ring-gray-100"
)
