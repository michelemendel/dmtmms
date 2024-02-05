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

// TODO: I think that some of this data belongs as regular args to the function, and not as part of the view context.
type Opts struct {
	UserOp               string
	Users                []entity.User
	SelectedUser         entity.User
	Roles                []string
	TempPassword         string
	TempPasswordUserName string
	Err                  error
	Msg                  string
}

func MakeOpts() Opts {
	return Opts{
		UserOp:               constants.OP_NONE,
		Users:                []entity.User{},
		SelectedUser:         entity.User{},
		Roles:                []string{"read", "edit", "admin"},
		TempPassword:         "",
		TempPasswordUserName: "",
		Err:                  nil,
		Msg:                  "",
	}
}

func (o Opts) WithOp(op string) Opts {
	o.UserOp = op
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

func (o Opts) WithMsg(msg string) Opts {
	o.Msg = msg
	return o
}
