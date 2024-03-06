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
	"github.com/michelemendel/dmtmms/util"
)

func (h *HandlerContext) DownloadHandler(c echo.Context) error {
	members, err := h.MembersFiltered(c)
	if err != nil {
		slog.Error("error getting members", "err", err)
		return err
	}

	var data []byte
	switch c.QueryParam("t") {
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
	// items = append(items, []string{"ID", "Name", "Fødselsnummer", "Email", "Mobile", "Address1", "Address2", "Postnummer", "Poststed", "Status", "RegisteredDate", "DeregisteredDate", "ReceiveEmail", "ReceiveMail", "ReceiveHatikvah", "Archived"})
	items = append(items, []string{"ID", "Name", "Family", "Fødselsnummer", "Age", "Email", "Mobile", "Address1", "Address2", "Postnummer", "Poststed", "Status", "RegisteredDate", "DeregisteredDate", "SynagogueSeat", "ReceiveEmail", "ReceiveMail", "ReceiveHatikvah", "Created", "Updated"})
	for _, m := range members {
		items = append(items, []string{
			util.Int2String(m.ID),
			m.Name,
			m.FamilyName,
			strings.Replace(util.Date2String(m.DOB), "-", "", -1) + "-" + m.Personnummer,
			fmt.Sprint(m.Age),
			m.Email,
			m.Mobile,
			m.Address1,
			m.Address2,
			m.Postnummer,
			m.Poststed,
			string(m.Status),
			util.Date2String(m.RegisteredDate),
			util.Date2String(m.DeregisteredDate),
			fmt.Sprint(m.Synagogueseat),
			util.Bool2String(m.ReceiveEmail),
			util.Bool2String(m.ReceiveMail),
			util.Bool2String(m.ReceiveHatikvah),
			// util.Bool2String(m.Archived),
			util.DateTime2String(m.CreatedAt),
			util.DateTime2String(m.UpdatedAt),
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
			strings.Replace(util.Date2String(m.DOB), "-", "", -1) + "-" + m.Personnummer,
		})
	}
	return items
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
