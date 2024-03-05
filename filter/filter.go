package filter

import (
	"fmt"
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
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
	SortCol        string
	SortOrder      string
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

func (o Opts) WSort(sort, order string) Opts {
	o.SortCol = sort
	o.SortOrder = order
	return o
}

//--------------------------------------------------------------------------------
// Filter from query parameters

func FilterFromQuery(c echo.Context) Filter {
	sort := c.QueryParam("sort")
	order := c.QueryParam("order")
	withSort := c.QueryParam("wsort")
	if sort == "" {
		sort = "f.name"
	}
	if withSort == "true" {
		if order == "" || order == "DESC" {
			order = "ASC"
		} else {
			order = "DESC"
		}
	}

	//
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
		WSelAges(selectedAges).
		WSort(sort, order),
	}
}

func (f Filter) URLQuery(mUUID, sortCol, sortOrder, withSort string) string {
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
		"&wsort=" + withSort +
		"&sort=" + sortCol +
		"&order=" + sortOrder +
		selAges
}

func (f Filter) URLForDownloadLink(downloadType string) string {
	return "/download" + f.URLQuery("", "", "", "false") + "&t=" + downloadType
}

func (f Filter) SortMembers(members []entity.Member) []entity.Member {
	fmt.Println("[Filter.SortMembers]:", f.SortCol, f.SortOrder)

	switch f.SortCol {
	case "ID":
		sort.Slice(members, func(i, j int) bool {
			if f.SortOrder == "ASC" {
				return members[i].ID < members[j].ID
			}
			return members[i].ID > members[j].ID
		})
	case "Name":
		sort.Slice(members, func(i, j int) bool {
			if f.SortOrder == "ASC" {
				return members[i].Name < members[j].Name
			}
			return members[i].Name > members[j].Name
		})
	case "Fødselsnummer":
		sort.Slice(members, func(i, j int) bool {
			if f.SortOrder == "ASC" {
				return members[i].DOB.Before(members[j].DOB)
			}
			return members[i].DOB.After(members[j].DOB)
		})
	case "Age":
		sort.Slice(members, func(i, j int) bool {
			if f.SortOrder == "ASC" {
				return members[i].Age < members[j].Age
			}
			return members[i].Age > members[j].Age
		})
	case "Family":
		sort.Slice(members, func(i, j int) bool {
			if f.SortOrder == "ASC" {
				return members[i].FamilyName < members[j].FamilyName
			}
			return members[i].FamilyName > members[j].FamilyName
		})
	case "RecEmail":
		sort.Slice(members, func(i, j int) bool {
			if f.SortOrder == "ASC" {
				return members[i].ReceiveEmail
			}
			return members[j].ReceiveEmail
		})
	case "RecMail":
		sort.Slice(members, func(i, j int) bool {
			if f.SortOrder == "ASC" {
				return members[i].ReceiveMail
			}
			return members[j].ReceiveMail
		})
	case "RecHatikvah":
		sort.Slice(members, func(i, j int) bool {
			if f.SortOrder == "ASC" {
				return members[i].ReceiveHatikvah
			}
			return members[j].ReceiveHatikvah
		})
	}

	return members
}
