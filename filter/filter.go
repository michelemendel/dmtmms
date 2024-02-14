package filter

import (
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
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
	MemberUUID  string
	FamilyUUID  string
	GroupUUID   string
	SearchTerms string
	From        string
	To          string
}

func MakeOpts() Opts {
	return Opts{
		GroupUUID: "",
		From:      constants.DATE_MIN,
		To:        constants.DATE_MAX,
	}
}

func (o Opts) WithMemberUUID(memberUUID string) Opts {
	o.MemberUUID = memberUUID
	return o
}

func (o Opts) WithFamilyUUID(familyUUID string) Opts {
	o.FamilyUUID = familyUUID
	return o
}

func (o Opts) WithGroupUUID(groupUUID string) Opts {
	o.GroupUUID = groupUUID
	return o
}

func (o Opts) WithSearchTerms(searchTerms string) Opts {
	o.SearchTerms = searchTerms
	return o
}

func (o Opts) WithFrom(from string) Opts {
	if from == "" {
		from = constants.DATE_MIN
	}
	o.From = from
	return o
}

func (o Opts) WithTo(to string) Opts {
	if to == "" {
		to = constants.DATE_MAX
	}
	o.To = to
	return o
}

//--------------------------------------------------------------------------------
// Filter from query parameters

func FilterFromQuery(c echo.Context) Filter {
	fuuid := c.QueryParam("fuuid")
	guuid := c.QueryParam("guuid")
	searchTerms := c.QueryParam("searchterms")
	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")
	return Filter{MakeOpts().WithFamilyUUID(fuuid).WithGroupUUID(guuid).WithSearchTerms(searchTerms).WithFrom(fromStr).WithTo(toStr)}
}
