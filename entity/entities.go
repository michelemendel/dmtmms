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
	UUID   string
	ID     string
	Name   string
	DOB    time.Time
	Email  string
	Mobile string
	Status MemberStatus
}

func NewMember(uuid, id, name string, dob time.Time, email, mobile string, status MemberStatus) Member {
	return Member{
		UUID:   uuid,
		ID:     id,
		Name:   name,
		DOB:    dob,
		Email:  email,
		Mobile: mobile,
		Status: status,
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
	Type string
}

func NewGroup(uuid, name, typ string) Group {
	return Group{
		UUID: uuid,
		Name: name,
		Type: typ,
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
	GType  GroupType
	MGRole string
}

func NewMemberGroupDTO(uuid, id, name string, dob time.Time, Email, Mobile string, status MemberStatus, guuid, gname string, gtype GroupType, mgRole string) MemberGroupDTO {
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
		GType:  gtype,
		MGRole: mgRole,
	}
}
