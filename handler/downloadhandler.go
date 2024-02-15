package handler

import (
	"bytes"
	"encoding/csv"
	"errors"
	"log/slog"

	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/filter"
	"github.com/michelemendel/dmtmms/util"
)

func (h *HandlerContext) MembersAsCSV(c echo.Context) ([]byte, error) {
	f := filter.FilterFromQuery(c)
	members, err := h.MembersFiltered(c, f)
	if err != nil {
		slog.Error("error getting members", "err", err)
		return nil, err

	}

	byteData, err := GetMembersAsByteArray(h.MembersTo2DArray(members))
	if err != nil {
		slog.Error("error creating csv data", "err", err)
		return nil, err
	}

	return byteData, nil
}

func (h *HandlerContext) MembersTo2DArray(members []entity.Member) [][]string {
	var items [][]string
	items = append(items, []string{"ID", "Name", "DOB", "Personnummer", "Email", "Mobile"})
	for _, m := range members {
		items = append(items, []string{
			m.ID,
			m.Name,
			util.Time2String(m.DOB),
			m.Personnummer,
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
