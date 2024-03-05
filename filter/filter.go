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
	SelectedAges    []string
	ReceiveEmail    string
	ReceiveMail     string
	ReceiveHatikvah string
	// Archived        string
	SelectedGroup  string
	SelectedStatus string
}

func MakeOpts() Opts {
	return Opts{
		GroupUUID: "",
		From:      constants.DATE_MIN,
		To:        constants.DATE_MAX,
	}
}

func (o Opts) WMemUUID(memberUUID string) Opts {
	o.MemberUUID = memberUUID
	return o
}

func (o Opts) WFamUUID(familyUUID string) Opts {
	o.FamilyUUID = familyUUID
	return o
}

func (o Opts) WGroupUUID(groupUUID string) Opts {
	o.GroupUUID = groupUUID
	return o
}

func (o Opts) WSearch(searchTerms string) Opts {
	o.SearchTerms = searchTerms
	return o
}

func (o Opts) WFrom(from string) Opts {
	// if from == "" {
	// from = constants.DATE_MIN
	// }
	o.From = from
	return o
}

func (o Opts) WTo(to string) Opts {
	// if to == "" {
	// to = constants.DATE_MAX
	// }
	o.To = to
	return o
}

func (o Opts) WSelAges(SelectedAges []string) Opts {
	o.SelectedAges = SelectedAges
	return o
}

func (o Opts) WRecEmail(receiveEmail string) Opts {
	o.ReceiveEmail = receiveEmail
	return o
}

func (o Opts) WRecMail(receiveMail string) Opts {
	o.ReceiveMail = receiveMail
	return o
}

func (o Opts) WRecHatikvah(receiveHatikvah string) Opts {
	o.ReceiveHatikvah = receiveHatikvah
	return o
}

// func (o Opts) WArchived(archived string) Opts {
// 	o.Archived = archived
// 	return o
// }

func (o Opts) WSelGroup(group string) Opts {
	o.SelectedGroup = group
	return o
}

func (o Opts) WSelStatus(status string) Opts {
	o.SelectedStatus = status
	return o
}

//--------------------------------------------------------------------------------
// Filter from query parameters

func FilterFromQuery(c echo.Context) Filter {
	muuid := c.QueryParam("muuid")
	fuuid := c.QueryParam("fuuid")
	guuid := c.QueryParam("guuid")
	searchTerms := c.QueryParam("searchterms")
	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")
	receiveEmail := c.QueryParam("receiveEmail")
	receiveMail := c.QueryParam("receiveMail")
	receiveHatikvah := c.QueryParam("receiveHatikvah")
	// archived := c.QueryParam("archived")
	selectedGroup := c.QueryParam("selectedGroup")
	selectedStatus := c.QueryParam("selectedStatus")
	selectedAges := c.QueryParams()["selectedAges"]

	return Filter{MakeOpts().
		WMemUUID(muuid).
		WFamUUID(fuuid).
		WGroupUUID(guuid).
		WSearch(searchTerms).
		WFrom(fromStr).
		WTo(toStr).
		WRecEmail(receiveEmail).
		WRecMail(receiveMail).
		WRecHatikvah(receiveHatikvah).
		// WArchived(archived).
		WSelGroup(selectedGroup).
		WSelStatus(selectedStatus).
		WSelAges(selectedAges),
	}
}

func (f Filter) URLQuery(mUUID string) string {
	selAges := ""
	for _, age := range f.SelectedAges {
		selAges += "&selectedAges=" + age
	}

	return "?muuid=" + mUUID +
		"&guuid=" + f.GroupUUID +
		"&from=" + f.From +
		"&to=" + f.To +
		"&searchterms=" + f.SearchTerms +
		"&selectedGroup=" + f.SelectedGroup +
		"&receiveEmail=" + f.ReceiveEmail +
		"&receiveMail=" + f.ReceiveMail +
		"&receiveHatikvah=" + f.ReceiveHatikvah +
		// "&archived=" + f.Archived +
		"&selectedStatus=" + f.SelectedStatus +
		selAges
}

func (f Filter) URLForDownloadLink(downloadType string) string {
	return "/download" + f.URLQuery("") + "&t=" + downloadType
}
