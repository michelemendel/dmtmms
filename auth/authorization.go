package auth

import "github.com/michelemendel/dmtmms/constants"

func IsAuthorized(loggedInUserRole string, object string) bool {
	switch loggedInUserRole {
	case "root":
		return true
	case "admin":
		switch object {
		case constants.AUTH_NAV_USERS:
			return true
		}
	case "read":
		switch object {
		case constants.AUTH_NAV_USERS:
			return false
		}
	}
	return false
}
