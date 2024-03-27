package filter

import (
	"sort"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
)

type Filter struct {
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
	SelectedGroup   string
	SelectedStatus  string
	SortCol         string
	SortOrder       string
	WithSort        string
}

func MakeFilterDefault() *Filter {
	return &Filter{
		GroupUUID: "",
		From:      constants.DATE_MIN,
		To:        constants.DATE_MAX,
	}
}

//--------------------------------------------------------------------------------
// Filter from query parameters

func (f *Filter) MakeFilterFromQuery(c echo.Context) {
	f.MemberUUID = c.QueryParam("muuid")
	f.FamilyUUID = c.QueryParam("fuuid")
	f.GroupUUID = c.QueryParam("guuid")
	f.From = c.QueryParam("from")
	f.To = c.QueryParam("to")
	f.SearchTerms = c.QueryParam("searchterms")
	f.ReceiveEmail = c.QueryParam("receiveEmail")
	f.ReceiveMail = c.QueryParam("receiveMail")
	f.ReceiveHatikvah = c.QueryParam("receiveHatikvah")
	f.SelectedGroup = c.QueryParam("selectedGroup")
	f.SelectedStatus = c.QueryParam("selectedStatus")
	f.SelectedAges = c.QueryParams()["selectedAges"]
}

func (f *Filter) URLForDownloadLink(downloadType string) string {
	return "/download" + f.URLQuery(f.SortCol, "false") + "&t=" + downloadType
}

func (f *Filter) URLQuery(sortCol, withSort string) string {
	selAges := ""
	for _, age := range f.SelectedAges {
		selAges += "&selectedAges=" + age
	}

	return "?fuuid=" + f.FamilyUUID +
		"&guuid=" + f.GroupUUID +
		"&from=" + f.From +
		"&to=" + f.To +
		"&searchterms=" + f.SearchTerms +
		"&selectedGroup=" + f.SelectedGroup +
		"&receiveEmail=" + f.ReceiveEmail +
		"&receiveMail=" + f.ReceiveMail +
		"&receiveHatikvah=" + f.ReceiveHatikvah +
		"&selectedStatus=" + f.SelectedStatus +
		"&wsort=" + withSort +
		"&sort=" + sortCol +
		"&order=" + f.SortOrder +
		selAges
}

func (f *Filter) SortMembers(c echo.Context, members []entity.Member) []entity.Member {
	col := c.QueryParam("sort")
	// order := c.QueryParam("order")
	withSort := c.QueryParam("wsort")

	// fmt.Println("   [SortMembers]:CURR:", f, f.SortCol, f.SortOrder, f.WithSort)
	// fmt.Println("   [SortMembers]:QUER:", f, col, order, withSort)

	if col == "" {
		f.SortCol = "Family"
		f.SortOrder = "ASC"
	}

	// f.SortCol = c.QueryParam("sort")
	// f.SortOrder = c.QueryParam("order")
	// f.WithSort = c.QueryParam("wsort")

	if withSort == "true" {
		f.WithSort = "false"
		if f.SortCol != col {
			f.SortOrder = "ASC"
			// fmt.Println("   [SortMembers]:SWAP:", f)
			f.SortCol = col
		} else if f.SortOrder == "ASC" {
			f.SortOrder = "DESC"
			// fmt.Println("   [SortMembers]:SWAP:", f)
		} else {
			f.SortOrder = "ASC"
			// fmt.Println("   [SortMembers]:SWAP:", f)
		}
	}

	// fmt.Println("   [SortMembers]: NEW:", f)

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
	case "Fødselsnr":
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
	case "SynSeat":
		sort.Slice(members, func(i, j int) bool {
			if f.SortOrder == "ASC" {
				return members[i].Synagogueseat < members[j].Synagogueseat
			}
			return members[i].Synagogueseat > members[j].Synagogueseat
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
