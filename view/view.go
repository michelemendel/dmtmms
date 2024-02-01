package view

import (
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
)

type ViewCtx struct {
	Opts
}

func MakeViewCtxDefault() *ViewCtx {
	return &ViewCtx{
		Opts: MakeOpts(),
	}
}

func MakeViewCtx(opts Opts) *ViewCtx {
	return &ViewCtx{
		Opts: opts,
	}
}

type Opts struct {
	Op                   string
	Users                []entity.User
	SelectedUser         entity.User
	Roles                []string
	TempPassword         string
	TempPasswordUserName string
	Err                  error
}

func MakeOpts() Opts {
	return Opts{
		Op:                   constants.OP_NONE,
		Users:                []entity.User{},
		SelectedUser:         entity.User{},
		Roles:                []string{"read", "admin"},
		TempPassword:         "",
		TempPasswordUserName: "",
		Err:                  nil,
	}
}

func (o Opts) WithOp(op string) Opts {
	o.Op = op
	return o
}

func (o Opts) WithUsers(users []entity.User) Opts {
	o.Users = users
	return o
}

func (o Opts) WithSelectedUser(user entity.User) Opts {
	o.SelectedUser = user
	return o
}

func (o Opts) WithRoles(roles []string) Opts {
	o.Roles = roles
	return o
}

func (o Opts) WithTempPW(password string, username string) Opts {
	o.TempPassword = password
	o.TempPasswordUserName = username
	return o
}

func (o Opts) WithErr(err error) Opts {
	o.Err = err
	return o
}