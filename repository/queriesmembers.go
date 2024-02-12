package repo

import (
	"database/sql"
	"log/slog"

	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
)

//--------------------------------------------------------------------------------
// Focus on members

const (
	_ = `
	SELECT m.uuid, m.id, m.name, date(m.dob), m.email, m.mobile, m.status
	FROM members as m 
	`

	queryMembersFamilies = `
	SELECT m.uuid, m.id, m.name, date(m.dob), m.email, m.mobile, m.status, f.name 
	FROM members as m 
	LEFT JOIN families as f ON m.family_uuid=f.uuid;
	`

	queryMembersGroups = `
	SELECT m.uuid, m.id, m.name, date(m.dob), m.email, m.mobile, m.status, IFNULL(f.uuid, ""), IFNULL(f.name, "")
	FROM members as m
	LEFT JOIN families as f ON m.family_uuid=f.uuid
	LEFT JOIN members_groups as mg on m.uuid = mg.member_uuid
	LEFT JOIN groups as g on mg.group_uuid = g.uuid
	WHERE m.dob BETWEEN julianday(?) AND julianday(?)
	`
)

func (r *Repo) SelectMembersByFilter(filter Filter) ([]entity.Member, error) {
	q := queryMembersGroups
	args := []any{
		filter.From,
		filter.To,
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
	var mUUID string
	var mId string
	var mName string
	var mDOBStr string
	var mEmail string
	var mMobile string
	var mStatus string
	var mFamilyUUID string
	var mFamilyName string
	for rows.Next() {
		err := rows.Scan(&mUUID, &mId, &mName, &mDOBStr, &mEmail, &mMobile, &mStatus, &mFamilyUUID, &mFamilyName)
		if err != nil {
			slog.Error(err.Error())
			return members, err
		}
		mDOB := util.String2Time(mDOBStr)
		members = append(members, entity.NewMember(mUUID, mId, mName, mDOB, mEmail, mMobile, entity.MemberStatus(mStatus), mFamilyUUID, mFamilyName))
	}
	err := rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return members, err
	}
	return members, nil
}

//--------------------------------------------------------------------------------
// Focus on groups

// Select all groups
func (r *Repo) SelectGroups() ([]entity.Group, error) {
	var groups []entity.Group
	var gUUID string
	var gName string
	var gType string
	rows, err := r.DB.Query("SELECT uuid, name, type FROM groups;")
	if err != nil {
		slog.Error(err.Error())
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&gUUID, &gName, &gType)
		if err != nil {
			slog.Error(err.Error())
			return groups, err
		}
		groups = append(groups, entity.NewGroup(gUUID, gName))
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return groups, err
	}
	return groups, nil
}

// Select groups by member
func (r *Repo) SelectGroupsByMember(memberUUID string) ([]entity.Group, error) {
	var groups []entity.Group
	var gUUID string
	var gName string
	rows, err := r.DB.Query(`
	SELECT g.uuid, g.name
	FROM groups as g 
	JOIN members_groups as mg on g.uuid = mg.group_uuid 
	WHERE mg.member_uuid = ?;
	`, memberUUID)
	if err != nil {
		slog.Error(err.Error())
		return groups, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&gUUID, &gName)
		if err != nil {
			slog.Error(err.Error())
			return groups, err
		}
		groups = append(groups, entity.NewGroup(gUUID, gName))
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return groups, err
	}
	return groups, nil
}

//--------------------------------------------------------------------------------
// Not sure we neeed this

// Select all members and groups
func (r *Repo) SelectMembersWithGroups() ([]entity.MemberGroupDTO, error) {
	var membersGroups []entity.MemberGroupDTO
	// Member
	var mUUID string
	var mId string
	var mName string
	var mDOB string
	var mEmail string
	var mMobile string
	var mStatus string
	// Group
	var gUUID string
	var gName string
	// Relationship
	var mgRole string

	rows, err := r.DB.Query(`
	SELECT m.uuid, m.id, m.name, m.dob, m.status, g.uuid, g.name, mg.role 
	FROM members as m 
	JOIN members_groups as mg on m.uuid = mg.member_uuid 
	JOIN groups as g on mg.group_uuid=g.uuid;
	`)
	if err != nil {
		slog.Error(err.Error())
		return membersGroups, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&mUUID, &mId, &mName, &mDOB, &mEmail, &mMobile, &mStatus)
		if err != nil {
			slog.Error(err.Error())
			return membersGroups, err
		}
		mDOB := util.String2Time(mDOB)
		membersGroups = append(membersGroups, entity.NewMemberGroupDTO(mUUID, mId, mName, mDOB, mEmail, mMobile, entity.MemberStatus(mStatus), gUUID, gName, mgRole))
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return membersGroups, err
	}
	return membersGroups, nil
}
