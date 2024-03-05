package entity

import (
	"time"
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
	ID           int
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
	ReceiveHatikvah   bool
	// Archived          bool
	Status     MemberStatus
	FamilyUUID string
	FamilyName string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewMember(
	uuid string,
	id int,
	name string,
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
	receiveHatikvah bool,
	// archived bool,
	status MemberStatus,
	familyUUID,
	familyName string,
	createdAt time.Time,
	updatedAt time.Time,
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
		ReceiveHatikvah:   receiveHatikvah,
		// Archived:          archived,
		Status:     status,
		FamilyUUID: familyUUID,
		FamilyName: familyName,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
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

type MemberDetails struct {
	MemberDetails []MemberDetail
	Groups        []Group
}
