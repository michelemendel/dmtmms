package repo

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
)

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
		mDOB, _ := time.Parse(constants.DATE_FRMT, mDOB)
		membersGroups = append(membersGroups, entity.NewMemberGroupDTO(mUUID, mId, mName, mDOB, mEmail, mMobile, entity.MemberStatus(mStatus), gUUID, gName, entity.GroupType(gType), mgRole))
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return membersGroups, err
	}
	return membersGroups, nil
}

//--------------------------------------------------------------------------------
// Focus on members

func (r *Repo) SelectMembers() ([]entity.Member, error) {
	query := `
	SELECT uuid, id, name, date(dob), email, mobile, status 
	FROM members;
	`
	return r.ExecuteQuery(query)
}

func (r *Repo) SelectMembersByGroup(groupUUID string) ([]entity.Member, error) {
	query := `
	SELECT m.uuid, m.id, m.name, m.dob, m.email, m.mobile, m.status 
	FROM members as m 
	JOIN members_groups as mg on m.uuid = mg.member_uuid 
	WHERE mg.group_uuid = ?;
	`
	return r.ExecuteQuery(query, groupUUID)
}

func (r *Repo) SelectMembersByDOBInterval(from time.Time, to time.Time) ([]entity.Member, error) {
	query := `
	SELECT uuid, id, name, dob, email, mobile, status 
	FROM members 
	WHERE dob BETWEEN julianday(?) AND julianday(?);
	`
	return r.ExecuteQuery(query, from, to)
}

func (r *Repo) ExecuteQuery(query string, args ...interface{}) ([]entity.Member, error) {
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
	var mUUID string
	var mId string
	var mName string
	var mDOB string
	var mEmail string
	var mMobile string
	var mStatus string
	for rows.Next() {
		err := rows.Scan(&mUUID, &mId, &mName, &mDOB, &mEmail, &mMobile, &mStatus)
		if err != nil {
			slog.Error(err.Error())
			return members, err
		}
		mDOB, _ := time.Parse(constants.DATE_FRMT, mDOB)
		members = append(members, entity.NewMember(mUUID, mId, mName, mDOB, mEmail, mMobile, entity.MemberStatus(mStatus)))
	}
	err := rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return members, err
	}
	return members, nil
}

// func (r *Repo) SelectMembers() ([]entity.Member, error) {
// 	var members []entity.Member
// 	var mUUID string
// 	var mId string
// 	var mName string
// 	var mDOB string
// 	var mEmail string
// 	var mMobile string
// 	var mStatus string
// 	rows, err := r.DB.Query("SELECT uuid, id, name, date(dob,'LOCALTIME'), status FROM members;")
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return members, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		err := rows.Scan(&mUUID, &mId, &mName, &mDOB, &mEmail, &mMobile, &mStatus)
// 		if err != nil {
// 			slog.Error(err.Error())
// 			return members, err
// 		}
// 		mDOB, _ := time.Parse(constants.DATE_FRMT, mDOB)
// 		members = append(members, entity.NewMember(mUUID, mId, mName, mDOB, mEmail, mMobile, entity.MemberStatus(mStatus)))
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return members, err
// 	}
// 	return members, nil
// }

// Select members by group
// func (r *Repo) SelectMembersByGroup(groupUUID string) ([]entity.Member, error) {
// 	var members []entity.Member
// 	var mUUID string
// 	var mId string
// 	var mName string
// 	var mDOB time.Time
// 	var mEmail string
// 	var mMobile string
// 	var mStatus string
// 	rows, err := r.DB.Query(`
// 	SELECT m.uuid, m.id, m.name, m.dob, m.status
// 	FROM members as m
// 	JOIN members_groups as mg on m.uuid = mg.member_uuid
// 	WHERE mg.group_uuid = ?;
// 	`, groupUUID)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return members, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		err := rows.Scan(&mUUID, &mId, &mName, &mDOB, &mEmail, &mMobile, &mStatus)
// 		if err != nil {
// 			slog.Error(err.Error())
// 			return members, err
// 		}
// 		members = append(members, entity.NewMember(mUUID, mId, mName, mDOB, mEmail, mMobile, entity.MemberStatus(mStatus)))
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return members, err
// 	}
// 	return members, nil
// }

// Select members by date of birth interval
// func (r *Repo) SelectMembersByDOBInterval(from time.Time, to time.Time) ([]entity.Member, error) {
// 	var members []entity.Member
// 	var mUUID string
// 	var mId string
// 	var mName string
// 	var mDOB time.Time
// 	var mEmail string
// 	var mMobile string
// 	var mStatus string
// 	rows, err := r.DB.Query(`
// 	SELECT uuid, id, name, dob, status
// 	FROM members
// 	WHERE dob BETWEEN julianday(?) AND julianday(?);
// 	`, from, to)
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return members, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		err := rows.Scan(&mUUID, &mId, &mName, &mDOB, &mEmail, &mMobile, &mStatus)
// 		if err != nil {
// 			slog.Error(err.Error())
// 			return members, err
// 		}
// 		members = append(members, entity.NewMember(mUUID, mId, mName, mDOB, mEmail, mMobile, entity.MemberStatus(mStatus)))
// 	}
// 	err = rows.Err()
// 	if err != nil {
// 		slog.Error(err.Error())
// 		return members, err
// 	}
// 	return members, nil
// }

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
		groups = append(groups, entity.NewGroup(gUUID, gName, gType))
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
		groups = append(groups, entity.NewGroup(gUUID, gName, gType))
	}
	err = rows.Err()
	if err != nil {
		slog.Error(err.Error())
		return groups, err
	}
	return groups, nil
}
