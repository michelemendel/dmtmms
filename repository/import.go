package repo

import (
	"encoding/csv"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path"
	"slices"
	"strconv"
	"time"

	"github.com/michelemendel/dmtmms/e"
	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
)

type Member struct {
	ID              int
	UUID            string
	Name            string
	DOB             time.Time
	Personnummer    string
	Mobile          string
	Address1        string
	Address2        string
	Postnummer      string
	Poststed        string
	Status          entity.MemberStatus
	FamilyUUID      string
	FamilyName      string
	DoesFamilyExist bool
}

func NewMember(id int, name string, dob time.Time, personnummer, mobile, address1, address2, postnummer, poststed, familyUUID, familyName string, doesFamilyExist bool) *Member {
	return &Member{
		ID:              id,
		UUID:            util.GenerateUUID(),
		Name:            name,
		DOB:             dob,
		Personnummer:    personnummer,
		Mobile:          mobile,
		Address1:        address1,
		Address2:        address2,
		Postnummer:      postnummer,
		Poststed:        poststed,
		Status:          entity.MemberStatusActive,
		FamilyUUID:      familyUUID,
		FamilyName:      familyName,
		DoesFamilyExist: doesFamilyExist,
	}
}

func (r *Repo) ImportData() {
	members, _ := GetMembers()

	for _, m := range members {
		if !m.DoesFamilyExist {
			f := entity.NewFamily(m.FamilyUUID, m.FamilyName)
			r.CreateFamily(f)
		}

		// fmt.Printf("id:%d, name:%s, dob:%s, personnummer:%s, mobile:%s, address1:%s, address2:%s, postnummer:%s, poststed:%s, familyName:%s\n", m.ID, m.Name, m.DOB, m.Personnummer, m.Mobile, m.Address1, m.Address2, m.Postnummer, m.Poststed, m.FamilyName)
		// fmt.Printf("id:%d, name:%s, dob:%s, dobType:%T, \n", m.ID, m.Name, m.DOB, m.DOB)

		r.ImportMember(m)
	}

	r.SetAutoIDTrigger()
}

func GetMembers() ([]Member, error) {
	fmt.Println("[import]")

	fileName := "members.csv"

	wd, _ := os.Getwd()
	fnFullPath := path.Join(wd, "cmd", "import", fileName)
	fmt.Println("Opening file", fnFullPath)
	file, err := os.Open(fnFullPath)
	if err != nil {
		fmt.Println("Error opening file", err)
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading file", err)
		return nil, err
	}

	var maxId = getMaxId(records)
	var ids []int
	var members []Member

	for _, record := range records {
		// CSV file header
		// Navn[0], Adresse-1[1], Adresse-2[2], Adresse-3[3], Telefon[4], F.dato[5], Personnummer[6], Medlem[7]
		idStr := record[7]
		name := record[0]
		dobStr := record[5]
		personnummerStr := record[6]
		telephone := record[4]
		address1 := record[1]
		address2 := record[2]
		address3 := record[3]

		//
		id, _ := strconv.Atoi(padId(idStr))
		personnummer := string(personnummerStr[6:])
		mobile := telephone[3:]
		postnummer := ""
		poststed := ""
		if len(record[3]) > 0 {
			postnummer = address3[0:4]
			poststed = address3[5:]
		}
		familyUUID := ""
		familyName := ""
		doesFamilyExist := false

		// fmt.Printf("%d - ", id)
		newID, isNewFamilyMember := maybeFixId(id, maxId, ids)
		ids = append(ids, id)
		if isNewFamilyMember {
			maxId = newID
			member, err := getFamilyMemberByID(id, members)
			if err == nil {
				// fmt.Printf("--- curr_member: name:%s, familyName:%s\n", member.Name, member.FamilyName)
				familyName = member.FamilyName
				familyUUID = member.FamilyUUID
			}
			doesFamilyExist = true
			id = newID
		} else {
			familyName = name
			familyUUID = util.GenerateUUID()
		}

		dob := dobStr[0:4] + "-" + dobStr[4:6] + "-" + dobStr[6:]
		m := NewMember(id, name, util.String2Date(dob), personnummer, mobile, address1, address2, postnummer, poststed, familyUUID, familyName, doesFamilyExist)
		members = append(members, *m)
	}

	return members, nil
}

func getFamilyMemberByID(id int, members []Member) (Member, error) {
	for _, m := range members {
		if m.ID == id {
			return m, nil
		}
	}
	return Member{}, errors.New("Member not found")
}

func maybeFixId(id, maxId int, ids []int) (int, bool) {
	if slices.Contains(ids, id) {
		newId := maxId + 1
		return newId, true
	}
	return id, false
}

func getMaxId(records [][]string) int {
	var ids []int
	for _, record := range records {
		idStr := record[7]
		id, _ := strconv.Atoi(padId(idStr))
		ids = append(ids, id)
	}
	return slices.Max(ids)
}

func padId(id string) string {
	length := len(id)
	switch length {
	case 1:
		id = "1000" + id
	case 2:
		id = "100" + id
	case 3:
		id = "10" + id
	case 4:
		id = "1" + id
	}
	return id
}

func (r *Repo) ImportMember(member Member) error {
	tx, _ := r.DB.Begin()

	_, err := tx.Exec(`
	INSERT INTO members(
		id,
		uuid, 
		name, dob, personnummer, mobile, 
		address1, address2, postnummer, poststed,  
		status,
		family_uuid, family_name
		) VALUES(
			?, 
			?, 
			?, julianday(?), ?, ?, 
			?, ?, ?, ?, 
			?,
			?, ?
		)
		`,
		member.ID,
		member.UUID,
		member.Name, member.DOB, member.Personnummer, member.Mobile,
		member.Address1, member.Address2, member.Postnummer, member.Poststed,
		member.Status,
		member.FamilyUUID, member.FamilyName,
	)
	if err != nil {
		slog.Error(err.Error(), "id", member.ID, "uuid", member.UUID, "name", member.Name)
		tx.Rollback()
		return e.ErrCreatingMember
	}

	tx.Commit()
	slog.Info("CreateMember", "id", member.ID, "uuid", member.UUID, "name", member.Name)
	return nil
}
