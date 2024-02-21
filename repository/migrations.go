package repo

import (
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
)

func (r *Repo) runStatements(sqlStmts map[string]string) {
	for name, sqlStmt := range sqlStmts {
		_, err := r.DB.Exec(sqlStmt)
		if err != nil {
			slog.Error(fmt.Sprintf("Error in stmt %q:%s\n%s\n", err, name, sqlStmt))
		}
	}
}

// --------------------------------------------------------------------------------
// DDL
func (r *Repo) DropTables() {
	var sqlStmts = make(map[string]string)

	sqlStmts["drop_users"] = `DROP TABLE IF EXISTS users;`
	sqlStmts["drop_members"] = `DROP TABLE IF EXISTS members;`
	sqlStmts["drop_families"] = `DROP TABLE IF EXISTS families;`
	sqlStmts["drop_groups"] = `DROP TABLE IF EXISTS groups;`
	sqlStmts["drop_members_groups"] = `DROP TABLE IF EXISTS members_groups;`

	r.runStatements(sqlStmts)
}

func (r *Repo) CreateTables() {
	var sqlStmts = make(map[string]string)

	sqlStmts["create_users"] = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP
	); `

	// id - this is an id shared by dmt and tripletex
	// https://mail.google.com/mail/u/0/#search/from%3Aadm%40dmt.oslo.no/FMfcgzGwJcbCTnSbdzxfhHXlFLZCXDGZ?projector=1&messagePartId=0.1
	// state = [active, archived, tobedeleted]
	sqlStmts["create_members"] = `
	CREATE TABLE IF NOT EXISTS members (
		uuid TEXT PRIMARY KEY,
		id TEXT,
		name TEXT NOT NULL,
		dob REAL,
		personnummer TEXT,
		email TEXT,
		mobile TEXT,
		address1 TEXT,
		address2 TEXT,
		postnummer TEXT,
		poststed TEXT,
		synagogue_seat TEXT,
		membership_fee_tier TEXT,
		registered_date REAL,
		deregistered_date REAL,
		receive_email BOOLEAN,
		receive_mail BOOLEAN,
		receive_hatikva BOOLEAN,
		status TEXT NOT NULL,
		archived BOOLEAN,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER,
		family_uuid TEXT,
		family_name TEXT,
		FOREIGN KEY(family_uuid) REFERENCES families(uuid)
	); `

	sqlStmts["create_families"] = `
	CREATE TABLE IF NOT EXISTS families (
		uuid TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER
	);`

	sqlStmts["create_groups"] = `
	CREATE TABLE IF NOT EXISTS groups (
		uuid TEXT PRIMARY KEY,
		name TEXT NOT NULL UNIQUE,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER
	);`

	// many-to-many between members and groups
	sqlStmts["create_members_groups"] = `
	CREATE TABLE IF NOT EXISTS members_groups (
		member_uuid TEXT NOT NULL,
		group_uuid TEXT NOT NULL,
		role TEXT,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER,
		primary key (member_uuid, group_uuid),
		FOREIGN KEY(member_uuid) REFERENCES members(uuid),
		FOREIGN KEY(group_uuid) REFERENCES groups(uuid)
	); `

	r.runStatements(sqlStmts)
}
func (r *Repo) CreateIndexes() {
	var sqlStmts = make(map[string]string)

	sqlStmts["create_index_members_uuid"] = `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_members_uuid ON members(uuid);`

	sqlStmts["create_index_families_uuid"] = `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_families_uuid ON families(uuid);`

	sqlStmts["create_index_groups_uuid"] = `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_groups_uuid ON groups(uuid);`

	sqlStmts["create_index_members_groups"] = `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_members_groups ON members_groups(member_uuid, group_uuid);`

	r.runStatements(sqlStmts)
}

//--------------------------------------------------------------------------------
// DML

// Users
func (r *Repo) InsertUsers() {
	users := []entity.User{
		{Name: "root", HashedPassword: "root", Role: "root"},
		{Name: "abe", HashedPassword: "abe", Role: "admin"},
		{Name: "eve", HashedPassword: "eve", Role: "edit"},
		{Name: "ron", HashedPassword: "ron", Role: "read"},
	}

	for _, user := range users {
		hpw, _ := util.HashPassword(user.HashedPassword)
		user.HashedPassword = hpw
		err := r.CreateUser(user)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

// Families
func (r *Repo) InsertFamilies() {
	families := []entity.Family{
		{UUID: "0", Name: ""},
		{UUID: "101", Name: "fam1"},
		{UUID: "102", Name: "fam2"},
	}

	for _, family := range families {
		err := r.CreateFamily(family)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

// Groups
func (r *Repo) InsertGroups() {
	groups := []entity.Group{
		{UUID: "0", Name: ""},
		{UUID: "1001", Name: "org1"},
		{UUID: "1002", Name: "org2"},
		{UUID: "1003", Name: "org3"},
		{UUID: "1004", Name: "org4"},
	}
	for _, group := range groups {
		err := r.CreateGroup(group)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

// Members
type member struct {
	uuid         string
	receiveEmail bool
	archived     bool
	familyUUID   string
	familyName   string
	groupUUIDs   []string
}

var members = map[string]member{
	"11": {uuid: "11", receiveEmail: true, archived: false, familyUUID: "101", familyName: "fam1", groupUUIDs: []string{"1001", "1002"}},
	"12": {uuid: "12", receiveEmail: true, archived: false, familyUUID: "101", familyName: "fam1", groupUUIDs: []string{"1002", "1003"}},
	"13": {uuid: "13", receiveEmail: false, archived: false, familyUUID: "102", familyName: "fam2", groupUUIDs: []string{"1001", "1002"}},
	"14": {uuid: "14", receiveEmail: true, archived: false, familyUUID: "102", familyName: "fam2", groupUUIDs: []string{"1002", "1003"}},
	"15": {uuid: "15", receiveEmail: false, archived: true, familyUUID: "102", familyName: "fam2", groupUUIDs: []string{"1002", "1004"}},
}

func (r *Repo) InsertMembers() {
	nofMembers := 7
	nrStart := 100
	memberUUID := 11
	dob := util.String2Time("1980-02-01")
	personnummer := "12345"
	memberIdPrefix := "99"
	namePrefix := "mem"
	phonePrefix := "12377"
	address1 := ""
	address2 := ""
	postnummer := ""
	poststed := ""
	var status entity.MemberStatus = entity.MemberStatusActive
	for i := 0; i < nofMembers; i++ {
		memberId := memberIdPrefix + strconv.Itoa(nrStart+i)
		name := namePrefix + strconv.Itoa(memberUUID)
		email := name + "@test.com"
		mobile := phonePrefix + strconv.Itoa(nrStart+i)
		address := entity.NewAddress(address1, address2, postnummer, poststed)
		receiveEmail := true
		archived := false
		familyUUID := ""
		familyName := ""
		groupUUIDs := []string{}
		if m, ok := members[strconv.Itoa(memberUUID)]; ok {
			receiveEmail = m.receiveEmail
			archived = m.archived
			familyUUID = m.familyUUID
			familyName = m.familyName
			groupUUIDs = m.groupUUIDs
		}

		member := entity.NewMember(strconv.Itoa(memberUUID), memberId, name, dob, personnummer, email, mobile, address, "", "", util.String2Time("2020-01-01"), time.Time{}, receiveEmail, false, false, archived, status, familyUUID, familyName)
		err := r.CreateMember(member, groupUUIDs)
		if err != nil {
			slog.Error(err.Error())
		}
		memberUUID++
		dob = dob.AddDate(0, 1, 1)
	}
}

func (r *Repo) ShowUsers() error {
	var name string
	var password string
	var role string
	var createdAt string
	// -- AND dob < julianday('1965-07-22')
	// rows, err := r.DB.Query(`SELECT name,date(dob),datetime(created_at,'LOCALTIME') FROM users;`)
	rows, err := r.DB.Query(`SELECT name,password,role,datetime(created_at,'LOCALTIME') FROM users;`)
	if err != nil {
		slog.Error("error getting users", "error", err.Error())
		return err
	}
	defer rows.Close()
	frmt := "%s\t%s\t%s\t%s\n"
	fmt.Printf(frmt, "name", "password", "role", "createdAt")
	for rows.Next() {
		err := rows.Scan(&name, &password, &role, &createdAt)
		if err != nil {
			slog.Error(err.Error())
		}
		fmt.Printf(frmt, name, password, role, createdAt)
	}

	return nil
}
