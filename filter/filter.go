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
	MemberUUID      string
	FamilyUUID      string
	GroupUUID       string
	SearchTerms     string
	From            string
	To              string
	ReceiveEmail    string
	ReceiveMail     string
	ReceiveHatikvah string
	Archived        string
	SelectedGroup   string
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
	// if from == "" {
	// from = constants.DATE_MIN
	// }
	o.From = from
	return o
}

func (o Opts) WithTo(to string) Opts {
	// if to == "" {
	// to = constants.DATE_MAX
	// }
	o.To = to
	return o
}

func (o Opts) WithReceiveEmail(receiveEmail string) Opts {
	o.ReceiveEmail = receiveEmail
	return o
}

func (o Opts) WithReceiveMail(receiveMail string) Opts {
	o.ReceiveMail = receiveMail
	return o
}

func (o Opts) WithReceiveHatikvah(receiveHatikvah string) Opts {
	o.ReceiveHatikvah = receiveHatikvah
	return o
}

func (o Opts) WithArchived(archived string) Opts {
	o.Archived = archived
	return o
}

func (o Opts) WithSelectedGroup(group string) Opts {
	o.SelectedGroup = group
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
	receiveEmail := c.QueryParam("receiveEmail")
	receiveMail := c.QueryParam("receiveMail")
	receiveHatikvah := c.QueryParam("receiveHatikvah")
	archived := c.QueryParam("archived")
	selectedGroup := c.QueryParam("selectedGroup")
	// fmt.Println("FilterFromQuery", fuuid, guuid, searchTerms, selectedGroup)
	return Filter{MakeOpts().WithFamilyUUID(fuuid).WithGroupUUID(guuid).WithSearchTerms(searchTerms).WithFrom(fromStr).WithTo(toStr).WithReceiveEmail(receiveEmail).WithReceiveMail(receiveMail).WithReceiveHatikvah(receiveHatikvah).WithArchived(archived).WithSelectedGroup(selectedGroup)}
}
