package repo

import (
	"database/sql"
	"log/slog"
	"strings"

	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/filter"
	"github.com/michelemendel/dmtmms/util"
)

//--------------------------------------------------------------------------------
// Members

const (
	selectMember = `
	SELECT 
	m.uuid, m.id, m.name, date(m.dob), m.personnummer, 
	m.email, m.mobile, 
	m.address1, m.address1, m.postnummer, m.poststed, 
	IFNULL(m.synagogue_seat, ""),
	IFNULL(m.membership_fee_tier, ""),
	IFNULL(date(m.registered_date), ""),
	IFNULL(date(m.deregistered_date), ""),
	IFNULL(m.receive_email, false),
	IFNULL(m.receive_mail, false),
	IFNULL(m.receive_hatikvah, false),
	IFNULL(m.archived, false),
	m.status, 
	IFNULL(f.uuid, ""), 
	IFNULL(f.name, "")
	FROM members as m
	`

	// queryMember = selectMember + `
	// LEFT JOIN families as f ON m.family_uuid=f.uuid
	// `

	queryMembers = selectMember + `
	LEFT JOIN families as f ON m.family_uuid=f.uuid
	LEFT JOIN members_groups as mg on m.uuid = mg.member_uuid
	LEFT JOIN groups as g on mg.group_uuid = g.uuid
	WHERE m.dob BETWEEN julianday(?) AND julianday(?)
	`
)

func (r *Repo) SelectMembersByFilter(filter filter.Filter) ([]entity.Member, error) {
	q := queryMembers
	from := constants.DATE_MIN
	to := constants.DATE_MAX
	if filter.From != "" {
		from = filter.From
	}
	if filter.To != "" {
		to = filter.To
	}
	args := []any{from, to}

	if strings.TrimSpace(filter.SearchTerms) != "" {
		q = q + "AND (m.name LIKE ? OR m.email LIKE ? OR f.name LIKE ?)"
		args = append(args, "%"+filter.SearchTerms+"%")
		args = append(args, "%"+filter.SearchTerms+"%")
		args = append(args, "%"+filter.SearchTerms+"%")
	}

	if filter.FamilyUUID != "" {
		q = q + "AND f.uuid=?"
		args = append(args, filter.FamilyUUID)
	}

	if filter.GroupUUID != "" {
		q = q + "AND g.uuid=?"
		args = append(args, filter.GroupUUID)
	}

	if filter.MemberUUID != "" {
		q = q + "AND m.uuid=?"
		args = append(args, filter.MemberUUID)
	}

	if filter.ReceiveEmail != "" {
		q = q + "AND m.receive_email=?"
		args = append(args, filter.ReceiveEmail)
	}

	if filter.ReceiveMail != "" {
		q = q + "AND m.receive_mail=?"
		args = append(args, filter.ReceiveMail)
	}

	if filter.ReceiveHatikvah != "" {
		q = q + "AND m.receive_hatikvah=?"
		args = append(args, filter.ReceiveHatikvah)
	}

	if filter.Archived != "" {
		q = q + "AND m.archived=?"
		args = append(args, filter.Archived)
	}

	if filter.SelectedGroup != "" && filter.SelectedGroup != "All groups" {
		q = q + "AND g.name=?"
		args = append(args, filter.SelectedGroup)
	}

	q = q + " GROUP BY m.uuid ORDER BY f.name;"
	return r.ExecuteQuery(q, args...)
}

//--------------------------------------------------------------------------------
//

func (r *Repo) ExecuteQuery(query string, args ...interface{}) ([]entity.Member, error) {
	// fmt.Println("ExecuteQuery", query, args)
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		slog.Error(err.Error())
		return []entity.Member{}, err
	}
	return r.MakeMemberList(rows)
}

// Run select on members
func (r *Repo) MakeMemberList(rows *sql.Rows) ([]entity.Member, error) {
	defer rows.Close()
	var members []entity.Member
	var dobStr string
	var registeredDateStr string
	var deregisteredDateStr string
	//
	var m entity.Member
	for rows.Next() {
		err := rows.Scan(
			&m.UUID, &m.ID, &m.Name, &dobStr, &m.Personnummer,
			&m.Email, &m.Mobile,
			&m.Address1, &m.Address2, &m.Postnummer, &m.Poststed,
			&m.Synagogueseat, &m.MembershipFeeTier, &registeredDateStr, &deregisteredDateStr,
			&m.ReceiveEmail, &m.ReceiveMail, &m.ReceiveHatikvah, &m.Archived, &m.Status, &m.FamilyUUID,
			&m.FamilyName,
		)
		if err != nil {
			slog.Error(err.Error())
			return members, err
		}

		DOB := util.String2Time(dobStr)
		registeredDate := util.String2Time(registeredDateStr)
		deregisteredDate := util.String2Time(deregisteredDateStr)
		address := entity.NewAddress(m.Address1, m.Address2, m.Postnummer, m.Poststed)
		members = append(members, entity.NewMember(
			m.UUID, m.ID, m.Name, DOB, m.Personnummer,
			m.Email, m.Mobile,
			address,
			m.Synagogueseat, m.MembershipFeeTier, registeredDate, deregisteredDate, m.ReceiveEmail,
			m.ReceiveMail, m.ReceiveHatikvah, m.Archived, entity.MemberStatus(m.Status), m.FamilyUUID,
			m.FamilyName,
		))
	}
	err := rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return members, err
	}
	return members, nil
}
