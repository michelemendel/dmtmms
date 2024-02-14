package repo

import (
	"database/sql"
	"log/slog"
	"strings"

	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/filter"
	"github.com/michelemendel/dmtmms/util"
)

//--------------------------------------------------------------------------------
// Focus on members

const (
	_ = `
	SELECT m.uuid, m.id, m.name, date(m.dob), m.email, m.mobile, m.status
	FROM members as m 
	`

	_ = `
	SELECT m.uuid, m.id, m.name, date(m.dob), m.email, m.mobile, m.status, f.name 
	FROM members as m 
	LEFT JOIN families as f ON m.family_uuid=f.uuid;
	`

	queryMembers = `
	SELECT m.uuid, m.id, m.name, date(m.dob), m.personnummer, m.email, m.mobile, m.address1, m.address1, m.postnummer, m.poststed, m.status, IFNULL(f.uuid, ""), IFNULL(f.name, "")
	FROM members as m
	LEFT JOIN families as f ON m.family_uuid=f.uuid
	LEFT JOIN members_groups as mg on m.uuid = mg.member_uuid
	LEFT JOIN groups as g on mg.group_uuid = g.uuid
	WHERE m.dob BETWEEN julianday(?) AND julianday(?)
	`
)

func (r *Repo) SelectMembersByFilter(filter filter.Filter) ([]entity.Member, error) {
	q := queryMembers
	args := []any{
		filter.From,
		filter.To,
	}

	// fmt.Printf("Searchterms: #%s#", filter.SearchTerms)
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

	q = q + " GROUP BY m.uuid ORDER BY m.name;"
	return r.ExecuteQuery(q, args...)
}

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
	var status string
	var familyUUID string
	var familyGroup string
	for rows.Next() {
		err := rows.Scan(&uuid, &id, &name, &dobStr, &personnummer, &email, &mobile, &address1, &address2, &postnummer, &poststed, &status, &familyUUID, &familyGroup)
		if err != nil {
			slog.Error(err.Error())
			return members, err
		}
		mDOB := util.String2Time(dobStr)
		address := entity.NewAddress(address1, address2, postnummer, poststed)
		members = append(members, entity.NewMember(uuid, id, name, mDOB, personnummer, email, mobile, address, entity.MemberStatus(status), familyUUID, familyGroup))
	}
	err := rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return members, err
	}
	return members, nil
}
