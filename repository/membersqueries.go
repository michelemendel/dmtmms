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
	m.status, 
	IFNULL(f.uuid, ""), IFNULL(f.name, ""),
	m.created_at, m.updated_at
	FROM members as m
	`

	queryMember = selectMember + `
	LEFT JOIN families as f ON m.family_uuid=f.uuid
	`

	queryMembers = selectMember + `
	LEFT JOIN families as f ON m.family_uuid=f.uuid
	LEFT JOIN members_groups as mg on m.uuid = mg.member_uuid
	LEFT JOIN groups as g on mg.group_uuid = g.uuid	
	`
)

func (r *Repo) SelectMembersByFilter(filter filter.Filter) ([]entity.Member, error) {
	q := queryMembers
	args := []any{}

	// Use either ages or from/to
	if len(filter.SelectedAges) > 0 && filter.SelectedAges[0] != "" {
		qs := []string{}
		for _, age := range filter.SelectedAges {
			year := util.GetYearFromAge(age)
			qs = append(qs, "strftime('%Y',dob)=? ")
			args = append(args, year)
		}
		q = q + "WHERE (" + strings.Join(qs, " OR ") + ") "
	} else {
		q = q + "WHERE m.dob BETWEEN julianday(?) AND julianday(?)"
		from := constants.DATE_MIN
		to := constants.DATE_MAX
		if filter.From != "" {
			from = filter.From
		}
		if filter.To != "" {
			to = filter.To
		}
		args = append(args, from)
		args = append(args, to)
	}

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

	// if filter.Archived != "" {
	// 	q = q + "AND m.archived=?"
	// 	args = append(args, filter.Archived)
	// }

	if filter.SelectedGroup != "" && filter.SelectedGroup != "All groups" {
		q = q + "AND g.name=?"
		args = append(args, filter.SelectedGroup)
	}

	if filter.SelectedStatus != "" {
		q = q + "AND m.status=?"
		args = append(args, filter.SelectedStatus)
	}

	q = q + " GROUP BY m.uuid ORDER BY f.name;"

	return r.ExecuteQuery(q, args...)
}

func (r *Repo) ExecuteQuery(query string, args ...interface{}) ([]entity.Member, error) {
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
	var createdAtStr string
	var updatedAtStr string
	//
	var m entity.Member
	for rows.Next() {
		err := rows.Scan(
			&m.UUID, &m.ID, &m.Name, &dobStr, &m.Personnummer,
			&m.Email, &m.Mobile,
			&m.Address1, &m.Address2, &m.Postnummer, &m.Poststed,
			&m.Synagogueseat, &m.MembershipFeeTier, &registeredDateStr, &deregisteredDateStr,
			// &m.ReceiveEmail, &m.ReceiveMail, &m.ReceiveHatikvah, &m.Archived, &m.Status, &m.FamilyUUID,
			&m.ReceiveEmail, &m.ReceiveMail, &m.ReceiveHatikvah, &m.Status, &m.FamilyUUID, &m.FamilyName,
			&createdAtStr, &updatedAtStr,
		)
		if err != nil {
			slog.Error(err.Error())
			return members, err
		}

		DOB := util.String2Date(dobStr)
		registeredDate := util.String2Date(registeredDateStr)
		deregisteredDate := util.String2Date(deregisteredDateStr)
		address := entity.NewAddress(m.Address1, m.Address2, m.Postnummer, m.Poststed)
		createdAt := util.String2DateTime(createdAtStr)
		updatedAt := util.String2DateTime(updatedAtStr)

		members = append(members, entity.NewMember(
			m.UUID, m.ID, m.Name, DOB, m.Personnummer,
			m.Email, m.Mobile,
			address,
			m.Synagogueseat, m.MembershipFeeTier, registeredDate, deregisteredDate, m.ReceiveEmail,
			// m.ReceiveMail, m.ReceiveHatikvah, m.Archived, entity.MemberStatus(m.Status), m.FamilyUUID,
			m.ReceiveMail, m.ReceiveHatikvah, entity.MemberStatus(m.Status), m.FamilyUUID,
			m.FamilyName, createdAt, updatedAt,
		))
	}
	err := rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return members, err
	}
	return members, nil
}

//--------------------------------------------------------------------------------
// Member

func (r *Repo) SelectMember(uuid string) (entity.Member, error) {
	q := queryMember + "WHERE m.uuid=?"
	args := []any{uuid}
	ms, err := r.ExecuteQuery(q, args...)
	if err != nil {
		return entity.Member{}, err
	}
	if len(ms) == 0 {
		return entity.Member{}, nil
	}

	return ms[0], nil
}
