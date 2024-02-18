package handler

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/filter"
	"github.com/michelemendel/dmtmms/util"
)

func (h *HandlerContext) DownloadHandler(c echo.Context) error {
	searchTerms := c.QueryParam("searchterms")
	from := c.QueryParam("from")
	to := c.QueryParam("to")
	downloadType := c.QueryParam("t")
	fmt.Println("searchTerms", searchTerms, "from", from, "to", to, "downloadType", downloadType)
	var data []byte
	members, _ := h.GetMembers(c)

	switch downloadType {
	case "csv":
		c.Response().Header().Set("Content-Disposition", "attachment; filename=members.csv")
		data, _ = h.MembersAsCSV(members, h.MembersTransformer)
		return c.Blob(200, "text/csv", data)
	case "fnr":
		c.Response().Header().Set("Content-Disposition", "attachment; filename=fodselsnummer.csv")
		data, _ = h.MembersAsCSV(members, h.PersonnummerTransformer)
		return c.Blob(200, "text/csv", data)
	case "emails":
		c.Response().Header().Set("Content-Disposition", "attachment; filename=emails.txt")
		data = h.Emails(members)
		return c.Blob(200, "text/plain", data)
	}
	return nil
}

func (h *HandlerContext) Emails(members []entity.Member) []byte {
	var emails string
	for _, m := range members {
		emails += m.Email + ","
	}
	return []byte(emails)
}

func (h *HandlerContext) MembersTransformer(members []entity.Member) [][]string {
	var items [][]string
	items = append(items, []string{"ID", "Name", "DOB", "Fødselsnummer", "Email", "Mobile"})
	for _, m := range members {
		items = append(items, []string{
			m.ID,
			m.Name,
			strings.Replace(util.Time2String(m.DOB), "-", "", -1) + "-" + m.Personnummer,
			m.Email,
			m.Mobile,
			// m.Address1,
			// m.Address2,
			// m.Postnummer,
			// m.Poststed,
			// m.Status,
		})
	}
	return items
}

func (h *HandlerContext) PersonnummerTransformer(members []entity.Member) [][]string {
	var items [][]string
	items = append(items, []string{"Namn", "Fødselsnummer"})
	for _, m := range members {
		items = append(items, []string{
			m.Name,
			strings.Replace(util.Time2String(m.DOB), "-", "", -1) + "-" + m.Personnummer,
		})
	}
	return items
}

func (h *HandlerContext) GetMembers(c echo.Context) ([]entity.Member, error) {
	f := filter.FilterFromQuery(c)
	members, err := h.MembersFiltered(c, f)
	if err != nil {
		slog.Error("error getting members", "err", err)
		return nil, err
	}
	return members, nil
}

func (h *HandlerContext) MembersAsCSV(members []entity.Member, transformFn func([]entity.Member) [][]string) ([]byte, error) {
	byteData, err := GetMembersAsByteArray(transformFn(members))
	if err != nil {
		slog.Error("error creating csv data", "err", err)
		return nil, err
	}

	return byteData, nil
}

func GetMembersAsByteArray(records [][]string) ([]byte, error) {
	if len(records) == 0 {
		return nil, errors.New("records cannot be nil or empty")
	}
	var buf bytes.Buffer
	csvWriter := csv.NewWriter(&buf)
	err := csvWriter.WriteAll(records)
	if err != nil {
		return nil, err
	}
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
