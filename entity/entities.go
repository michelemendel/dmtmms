package entity

import (
	"time"

	"github.com/michelemendel/dmtmms/util"
)

type InputError struct {
	InputID string // ex: "name", "email", "mobile"
	Err     error
}

func NewInputError(inputID string, err error) InputError {
	return InputError{
		InputID: inputID,
		Err:     err,
	}
}

type InputErrors map[string]InputError

func NewInputErrors() InputErrors {
	return InputErrors{}
}

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

// Aktiv, Utmeldt, Død, Stopp
// TODO: What does "Stopp" mean?
const (
	MemberStatusActive       MemberStatus = "active"
	MemberStatusDeregistered MemberStatus = "deregistered"
	MemberStatusDead         MemberStatus = "dead"
	MemberStatusStop         MemberStatus = "stop"
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
// ok Navn
// ok Adresse
// ok Kommune
// ok Fødselsdato
// ok Personnummer
// ok E-post
// ok Telefonnummer
// ok Synagogeplass
// ok Medlemsbidrags-grupper
// ok Innmeldings dato
// ok Utmeldings dato
// ok Status medlemskap Aktiv, Utmeldt, Død, Stopp
// ok Hatikva?
type Member struct {
	UUID         string
	ID           string
	Name         string
	DOB          time.Time
	Personnummer string
	Email        string
	Mobile       string
	Address
	Synagogueseat     string
	MembershipFeeTier string
	RegisteredDate    time.Time
	DeregisteredDate  time.Time
	ReceiveEmail      bool
	ReceiveMail       bool
	ReceiveHatikva    bool
	Archived          bool
	Status            MemberStatus
	FamilyUUID        string
	FamilyName        string
}

func NewMember(uuid,
	id, name string,
	dob time.Time,
	personnummer,
	email,
	mobile string,
	address Address,
	synagogueseat string,
	membershipFeeTier string,
	registeredDate time.Time,
	deregisteredDate time.Time,
	receiveEmail bool,
	receiveMail bool,
	receiveHatikva bool,
	archived bool,
	status MemberStatus,
	familyUUID,
	familyName string,
) Member {
	return Member{
		UUID:              uuid,
		ID:                id,
		Name:              name,
		DOB:               dob,
		Personnummer:      personnummer,
		Email:             email,
		Mobile:            mobile,
		Address:           address,
		Synagogueseat:     synagogueseat,
		MembershipFeeTier: membershipFeeTier,
		RegisteredDate:    registeredDate,
		DeregisteredDate:  deregisteredDate,
		ReceiveEmail:      receiveEmail,
		ReceiveMail:       receiveMail,
		ReceiveHatikva:    receiveHatikva,
		Archived:          archived,
		Status:            status,
		FamilyUUID:        familyUUID,
		FamilyName:        familyName,
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
// Maybe not the best way to get a list of data in order, when we have to write memberDetails[2].Value to get some value.
func GetMemberDetailsForPresentation(member Member) []MemberDetail {
	details := []MemberDetail{}

	if member.UUID == "" {
		return details
	}

	details = append(details, MemberDetail{"UUID", member.UUID})
	details = append(details, MemberDetail{"FamilyUUID", member.FamilyUUID})
	details = append(details, MemberDetail{"FamilyName", member.FamilyName})
	details = append(details, MemberDetail{"Name", member.Name})
	details = append(details, MemberDetail{"ID", member.ID})
	details = append(details, MemberDetail{"Date of Birth", util.Time2String(member.DOB)})
	details = append(details, MemberDetail{"Personnummer", member.Personnummer})
	details = append(details, MemberDetail{"Email", member.Email})
	details = append(details, MemberDetail{"Mobile", member.Mobile})
	details = append(details, MemberDetail{"Address1", member.Address.Address1})
	details = append(details, MemberDetail{"Address2", member.Address.Address2})
	details = append(details, MemberDetail{"Poststed", member.Address.Postnummer + " " + member.Address.Poststed})
	details = append(details, MemberDetail{"Status", string(member.Status)})
	details = append(details, MemberDetail{"Synagogueseat", member.Synagogueseat})
	details = append(details, MemberDetail{"MembershipFeeTier", member.MembershipFeeTier})
	details = append(details, MemberDetail{"RegisteredDate", util.Time2String(member.RegisteredDate)})
	details = append(details, MemberDetail{"DeregisteredDate", util.Time2String(member.DeregisteredDate)})
	details = append(details, MemberDetail{"ReceiveEmail", util.Bool2String(member.ReceiveEmail)})
	details = append(details, MemberDetail{"ReceiveMail", util.Bool2String(member.ReceiveMail)})
	details = append(details, MemberDetail{"ReceiveHatikva", util.Bool2String(member.ReceiveHatikva)})
	details = append(details, MemberDetail{"Archived", util.Bool2String(member.Archived)})

	return details
}
