package constants

const (
	DATE_FRMT = "2006-01-02"
	DATE_MIN  = "1000-01-01"
	DATE_MAX  = "3000-01-01"
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
)

// Routes
const (
	//
	ROUTE_LOGIN  = "/login"
	ROUTE_LOGOUT = "/logout"
	//
	ROUTE_INDEX = "/"
	//
	// Users
	ROUTE_USERS           = "/users"
	ROUTE_USER_CREATE     = "/user/create"
	ROUTE_USER_DELETE     = "/user/delete"
	ROUTE_USER_UPDATE     = "/user/update"
	ROUTE_USER_RESET_PW   = "/user/resetpw"
	ROUTE_USER_SET_PW     = "/user/setpw"
	ROUTE_USER_EDIT_CLOSE = "/user/editclose"
	//
	// Members
	ROUTE_MEMBERS     = "/members"
	ROUTE_MEMBER_EDIT = "/member/edit"
	//
	// Groups
	ROUTE_GROUPS           = "/groups"
	ROUTE_GROUP_CREATE     = "/group/create"
	ROUTE_GROUP_DELETE     = "/group/delete"
	ROUTE_GROUP_UPDATE     = "/group/update"
	ROUTE_GROUP_EDIT_CLOSE = "/group/editclose"
	//
	// Families
	ROUTE_FAMILIES          = "/families"
	ROUTE_FAMILY_CREATE     = "/family/create"
	ROUTE_FAMILY_DELETE     = "/family/delete"
	ROUTE_FAMILY_UPDATE     = "/family/update"
	ROUTE_FAMILY_EDIT_CLOSE = "/family/editclose"
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
