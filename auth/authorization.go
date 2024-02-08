package auth

import (
	"github.com/michelemendel/dmtmms/constants"
)

func IsAuthorized(loggedInUserRole string, target string) bool {
	switch loggedInUserRole {
	case "root":
		return true
	case "admin":
		switch target {
		case constants.AUTH_NAV_USERS:
			return true
		}
	case "edit":
		switch target {
		case constants.AUTH_NAV_USERS:
			return false
		}
	case "read":
		switch target {
		case constants.AUTH_NAV_USERS:
			return false
		}
	}
	return false
}
