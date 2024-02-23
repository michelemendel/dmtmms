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
		q = q + "AND (m.name LIKE ? OR m.email LIKE ?)"
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

	q = q + " GROUP BY m.uuid ORDER BY m.name;"
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
	//
	var uuid string
	var id string
	var name string
	var dobStr string
	var personnummer string
	var email string
	var mobile string
	var address1 string
	var address2 string
	var postnummer string
	var poststed string
	var synagogueSeat string
	var membershipFeeTier string
	var registeredDateStr string
	var deregisteredDateStr string
	var receiveEmail bool
	var receiveMail bool
	var receiveHatikva bool
	var archived bool
	var status string
	var familyUUID string
	var familyName string
	//
	for rows.Next() {
		err := rows.Scan(
			&uuid, &id, &name, &dobStr, &personnummer,
			&email, &mobile,
			&address1, &address2, &postnummer, &poststed,
			&synagogueSeat, &membershipFeeTier, &registeredDateStr, &deregisteredDateStr,
			&receiveEmail, &receiveMail, &receiveHatikva, &archived, &status, &familyUUID,
			&familyName,
		)
		if err != nil {
			slog.Error(err.Error())
			return members, err
		}

		DOB := util.String2Time(dobStr)
		registeredDate := util.String2Time(registeredDateStr)
		deregisteredDate := util.String2Time(deregisteredDateStr)
		address := entity.NewAddress(address1, address2, postnummer, poststed)
		members = append(members, entity.NewMember(
			uuid, id, name, DOB, personnummer,
			email, mobile,
			address,
			synagogueSeat, membershipFeeTier, registeredDate, deregisteredDate, receiveEmail,
			receiveMail, receiveHatikva, archived, entity.MemberStatus(status), familyUUID,
			familyName,
		))
	}
	err := rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return members, err
	}
	return members, nil
}
