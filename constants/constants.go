package constants

// Env keys
const (
	APP_ENV_KEY         = "APP_ENV"
	LOG_OUTPUT_TYPE_KEY = "LOG_OUTPUT_TYPE"
	LOG_FILE_NAME_KEY   = "LOG_FILE_NAME"
	WEB_SERVER_PORT_KEY = "WEB_SERVER_PORT"
	SESSION_KEY_KEY     = "SESSION_KEY"
	DEV_DB_DIR_KEY      = "DEV_DB_DIR"
	PROD_DB_DIR_KEY     = "PROD_DB_DIR"
	DB_NAME_KEY         = "DB_NAME"
)

const (
	LOG_OUTPUT_TYPE_FILE   = "file"
	LOG_OUTPUT_TYPE_STDOUT = "stdout"
	SESSION_NAME           = "session"
	TOKEN_NAME             = "token"
)

// Context keys
const (
	ERROR_KEY     = "error"
	USER_NAME_KEY = "username"
	// TokenKey      = "token"
	IS_LOGGEDIN_KEY = "isloggedinkey"
)
