package view

import "github.com/michelemendel/dmtmms/auth"

type ViewCtx struct {
	Opts
	Session *auth.Session
}

func MakeViewCtxDefault(session *auth.Session) *ViewCtx {
	return &ViewCtx{
		Opts:    MakeOpts(),
		Session: session,
	}
}

func MakeViewCtx(session *auth.Session, opts Opts) *ViewCtx {
	return &ViewCtx{
		Opts:    opts,
		Session: session,
	}
}

type ErrType string

const (
	ErrTypeNone     ErrType = ""
	ErrTypeOnCreate ErrType = "onCreate"
	ErrTypeOnUpdate ErrType = "onUpdate"
	ErrTypeOnDelete ErrType = "onDelete"
)

type ViewError struct {
	ErrType ErrType
	Err     error
}

type Opts struct {
	Roles                []string
	TempPassword         string
	TempPasswordUserName string
	ViewError            ViewError
	Msg                  string
}

func MakeOpts() Opts {
	return Opts{
		Roles:                []string{"read", "edit", "admin"},
		TempPassword:         "",
		TempPasswordUserName: "",
		ViewError:            ViewError{},
		Msg:                  "",
	}
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
	o.ViewError = ViewError{Err: err}
	return o
}

func (o Opts) WithErrType(err error, errType ErrType) Opts {
	o.ViewError = ViewError{Err: err, ErrType: errType}
	return o
}

func (o Opts) WithMsg(msg string) Opts {
	o.Msg = msg
	return o
}
