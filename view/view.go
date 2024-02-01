package view

import "github.com/michelemendel/dmtmms/entity"

type ViewCtx struct {
	Op           string
	Users        []entity.User
	SelectedUser entity.User
	Roles        []string
	Err          error
}

func MakeViewCtx(users []entity.User, selectedUser entity.User, op string, err error) *ViewCtx {
	return &ViewCtx{
		Op:           op,
		Users:        users,
		SelectedUser: selectedUser,
		Roles:        []string{"read", "admin"},
		Err:          err,
	}
}
