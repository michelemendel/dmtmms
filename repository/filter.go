package repo

import "time"

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
	From      time.Time
	To        time.Time
	GroupUUID string
}

func MakeOpts() Opts {
	return Opts{
		From:      time.Time{},
		To:        time.Time{},
		GroupUUID: "",
	}
}

func (o Opts) WithFrom(from time.Time) Opts {
	o.From = from
	return o
}

func (o Opts) WithTo(to time.Time) Opts {
	o.To = to
	return o
}

func (o Opts) WithGroupUUID(groupUUID string) Opts {
	o.GroupUUID = groupUUID
	return o
}
