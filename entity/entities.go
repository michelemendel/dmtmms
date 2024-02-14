package entity

import (
	"time"

	"github.com/michelemendel/dmtmms/util"
)

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

// https://www.posten.no/sende/adressering#:~:text=Adressering%20til%20postmottakere%20i%20Norge&text=Felles%20for%20all%20adressering%20er,p%C3%A5%20nederste%20linje%20i%20adressen.
type Address struct {
	Address1   string
	Address2   string
	Postnummer string
	Poststed   string
}

func NewAddress(address1, address2, postnummer, poststed string) Address {
	return Address{
		Address1:   address1,
		Address2:   address2,
		Postnummer: postnummer,
		Poststed:   poststed,
	}
}

// Fødselsnummer (11 digits) = dateOfBirth (6 digits) + personnummer (5 digits)
// https://mail.google.com/mail/u/0/#search/from%3Aadm%40dmt.oslo.no/FMfcgzGwJcbCTnSbdzxfhHXlFLZCXDGZ?projector=1&messagePartId=0.1
type Member struct {
	UUID         string
	ID           string
	Name         string
	DOB          time.Time
	Personnummer string
	Email        string
	Mobile       string
	Status       MemberStatus
	Address
	ReceiveEmail bool
	ReceiveMail  bool
	FamilyUUID   string
	FamilyGroup  string
}

func NewMember(uuid,
	id, name string,
	dob time.Time,
	personnummer,
	email, mobile string,
	address Address,
	status MemberStatus,
	familyUUID, familyGroup string,
) Member {
	return Member{
		UUID:        uuid,
		ID:          id,
		Name:        name,
		DOB:         dob,
		Email:       email,
		Mobile:      mobile,
		Address:     address,
		Status:      status,
		FamilyUUID:  familyUUID,
		FamilyGroup: familyGroup,
	}
}

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

type Family struct {
	UUID string
	Name string
}

func NewFamily(uuid, name string) Family {
	return Family{
		UUID: uuid,
		Name: name,
	}
}

//--------------------------------------------------------------------------------
// MemberDatas

type MemberDetail struct {
	Title string
	Value string
}

// To be presented on the detail section of the member page
func GetMemberDetails(member Member) []MemberDetail {
	datas := []MemberDetail{}

	if member.DOB.IsZero() {
		return datas
	}

	datas = append(datas, MemberDetail{"ID", member.ID})
	datas = append(datas, MemberDetail{"Date of Birth", util.Time2String(member.DOB)})
	datas = append(datas, MemberDetail{"Personnummer", member.Personnummer})
	datas = append(datas, MemberDetail{"Name", member.Name})
	datas = append(datas, MemberDetail{"Email", member.Email})
	datas = append(datas, MemberDetail{"Mobile", member.Mobile})
	datas = append(datas, MemberDetail{"Address1", member.Address.Address1})
	datas = append(datas, MemberDetail{"Address2", member.Address.Address2})
	datas = append(datas, MemberDetail{"Poststed", member.Address.Postnummer + " " + member.Address.Poststed})
	datas = append(datas, MemberDetail{"Status", string(member.Status)})

	return datas
}
