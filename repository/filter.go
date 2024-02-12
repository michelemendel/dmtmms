package repo

import (
	"time"

	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/util"
)

type Filter struct {
	Opts
}

func MakeFilterDefault() *Filter {
	return &Filter{
		Opts: MakeOpts(),
	}
}

func MakeFilter(opts Opts) *Filter {
	return &Filter{
		Opts: opts,
	}
}

type Opts struct {
	From       time.Time
	To         time.Time
	GroupUUID  string
	MemberUUID string
}

func MakeOpts() Opts {
	return Opts{
		From:      util.String2Time(constants.DATE_MIN),
		To:        util.String2Time(constants.DATE_MAX),
		GroupUUID: "",
	}
}

func (o Opts) WithFrom(fromStr string) Opts {
	if fromStr == "" {
		fromStr = constants.DATE_MIN
	}
	o.From = util.String2Time(fromStr)
	return o
}

func (o Opts) WithTo(toStr string) Opts {
	if toStr == "" {
		toStr = constants.DATE_MAX
	}
	o.To = util.String2Time(toStr)
	return o
}

func (o Opts) WithGroupUUID(groupUUID string) Opts {
	o.GroupUUID = groupUUID
	return o
}

func (o Opts) WithMemberUUID(memberUUID string) Opts {
	o.MemberUUID = memberUUID
	return o
}
