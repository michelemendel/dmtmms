package repo

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
)

//--------------------------------------------------------------------------------
// Focus on members

const (
	queryMembers = `
	SELECT m.uuid, m.id, m.name, date(m.dob), m.email, m.mobile, m.status, g.name, g.type
	FROM members as m 
	JOIN members_groups as mg on m.uuid = mg.member_uuid 
	JOIN groups as g on mg.group_uuid = g.uuid	
	GROUP BY m.uuid
	`
)

func (r *Repo) SelectMembersByFilter(filter Filter) ([]entity.Member, error) {
	if filter.GroupUUID != "" {
		return r.SelectMembersByGroupUUID(filter.GroupUUID)
	}

	// TODO: Fix, since this is always true
	if !filter.From.IsZero() && !filter.To.IsZero() {
		return r.SelectMembersByDOBInterval(filter.From, filter.To)
	}

	return r.SelectMembers()
}

func (r *Repo) SelectMembers() ([]entity.Member, error) {
	return r.ExecuteQuery(queryMembers + ";")
}

func (r *Repo) SelectMembersByGroupUUID(groupUUID string) ([]entity.Member, error) {
	query := queryMembers + " WHERE g.uuid = ?;"
	return r.ExecuteQuery(query, groupUUID)
}

func (r *Repo) SelectMembersByDOBInterval(from time.Time, to time.Time) ([]entity.Member, error) {
	query := queryMembers + " WHERE m.dob BETWEEN julianday(?) AND julianday(?);"
	return r.ExecuteQuery(query, util.Time2String(from), util.Time2String(to))
}

func (r *Repo) ExecuteQuery(query string, args ...interface{}) ([]entity.Member, error) {
	// fmt.Println("ExecuteQuery", query, args)
	rows, err := r.DB.Query(query, args...)
	if err != nil {
		slog.Error(err.Error())
		return []entity.Member{}, err
	}
	return r.SQLMembers(rows)
}

// Run select on members
func (r *Repo) SQLMembers(rows *sql.Rows) ([]entity.Member, error) {
	defer rows.Close()
	var members []entity.Member
	//
	var mUUID string
	var mId string
	var mName string
	var mDOB string
	var mEmail string
	var mMobile string
	var mStatus string
	var gGroupName string
	var gGroupType string
	for rows.Next() {
		err := rows.Scan(&mUUID, &mId, &mName, &mDOB, &mEmail, &mMobile, &mStatus, &gGroupName, &gGroupType)
		if err != nil {
			slog.Error(err.Error())
			return members, err
		}
		mDOB := util.String2Time(mDOB)
		members = append(members, entity.NewMember(mUUID, mId, mName, mDOB, mEmail, mMobile, entity.MemberStatus(mStatus), gGroupName, gGroupType))
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

func (r *Repo) SelectMemberByUUID(memberUUID string) (entity.Member, error) {
	query := queryMembers + " AND m.uuid = ?;"

	var mUUID string
	var mId string
	var mName string
	var mDOBStr string
	var mEmail string
	var mMobile string
	var mStatus string
	var gGroupName string
	var gGroupType string
	err := r.DB.QueryRow(query, memberUUID).Scan(&mUUID, &mId, &mName, &mDOBStr, &mEmail, &mMobile, &mStatus, &gGroupName, &gGroupType)
	if err != nil {
		slog.Error("Couldn't find member", "uuid", memberUUID, "error", err.Error())
		return entity.Member{}, err
	}

	mDOB := util.String2Time(mDOBStr)
	return entity.NewMember(mUUID, mId, mName, mDOB, mEmail, mMobile, entity.MemberStatus(mStatus), gGroupName, gGroupType), nil
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
		groups = append(groups, entity.NewGroup(gUUID, gName, entity.GroupType(gType)))
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
	var gType string
	rows, err := r.DB.Query(`
	SELECT g.uuid, g.name, g.type 
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
		err := rows.Scan(&gUUID, &gName, &gType)
		if err != nil {
			slog.Error(err.Error())
			return groups, err
		}
		groups = append(groups, entity.NewGroup(gUUID, gName, entity.GroupType(gType)))
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
	var gType string
	// Relationship
	var mgRole string

	rows, err := r.DB.Query(`
	SELECT m.uuid, m.id, m.name, m.dob, m.status, g.uuid, g.name, g.type, mg.role 
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
		membersGroups = append(membersGroups, entity.NewMemberGroupDTO(mUUID, mId, mName, mDOB, mEmail, mMobile, entity.MemberStatus(mStatus), gUUID, gName, entity.GroupType(gType), mgRole))
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return membersGroups, err
	}
	return membersGroups, nil
}
