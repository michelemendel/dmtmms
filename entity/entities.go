package entity

import "time"

// DTO
type User struct {
	Name           string
	HashedPassword string
	Role           string
}

func NewUser(name, password, role string) User {
	return User{
		Name:           name,
		HashedPassword: password,
		Role:           role,
	}
}

type MemberStatus string

const (
	MemberStatusActive      MemberStatus = "active"
	MemberStatusArchived    MemberStatus = "archived"
	MemberStatusToBeDeleted MemberStatus = "tobedeleted"
)

type Member struct {
	UUID       string
	ID         string
	Name       string
	DOB        time.Time
	Email      string
	Mobile     string
	Status     MemberStatus
	FamilyUUID string
	FamilyName string
}

func NewMember(uuid, id, name string, dob time.Time, email, mobile string, status MemberStatus, familyUUID, familyName string) Member {
	return Member{
		UUID:       uuid,
		ID:         id,
		Name:       name,
		DOB:        dob,
		Email:      email,
		Mobile:     mobile,
		Status:     status,
		FamilyUUID: familyUUID,
		FamilyName: familyName,
	}
}

type GroupType string

const (
	GroupTypeFamily GroupType = "fam"
	GroupTypeOrg    GroupType = "org"
)

type Group struct {
	UUID string
	Name string
}

func NewGroup(uuid, name string) Group {
	return Group{
		UUID: uuid,
		Name: name,
	}
}

type MemberGroupDTO struct {
	UUID   string
	ID     string
	Name   string
	DOB    time.Time
	Email  string
	Mobile string
	Status MemberStatus
	GUUID  string
	GName  string
	MGRole string
}

func NewMemberGroupDTO(uuid, id, name string, dob time.Time, Email, Mobile string, status MemberStatus, guuid, gname string, mgRole string) MemberGroupDTO {
	return MemberGroupDTO{
		UUID:   uuid,
		ID:     id,
		Name:   name,
		DOB:    dob,
		Email:  Email,
		Mobile: Mobile,
		Status: status,
		GUUID:  guuid,
		GName:  gname,
		MGRole: mgRole,
	}
}
