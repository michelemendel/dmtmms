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
		id TEXT UNIQUE,
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
		family_group TEXT,
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

type user struct {
	name string
	pw   string
	role string
}

func (r *Repo) InsertUsers() {
	stmt, err := r.DB.Prepare("INSERT INTO users(name,password,role) values(?, ?, ?)")
	if err != nil {
		slog.Error(err.Error())
	}

	users := []user{
		{"root", "root", "root"},
		{"abe", "abe", "admin"},
		{"eve", "eve", "edit"},
		{"ron", "ron", "read"},
	}

	for _, user := range users {
		hpw, _ := util.HashPassword(user.pw)
		_, err = stmt.Exec(user.name, hpw, user.role)
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

// Members

func (r *Repo) InsertFamilies() {
	familyStmt, _ := r.DB.Prepare("INSERT INTO families(uuid, name) values(?, ?)")
	familyUUID := 101
	families := []string{"fam1", "fam2"}
	for _, familyGroup := range families {
		_, err := familyStmt.Exec(familyUUID, familyGroup)
		if err != nil {
			slog.Error(err.Error())
		}
		familyUUID++
	}
}
func (r *Repo) InsertGroups() {
	groupStmt, _ := r.DB.Prepare("INSERT INTO groups(uuid, name) values(?, ?)")
	groupUUID := 1001
	groups := []string{"org1", "org2", "org3"}
	for _, groupName := range groups {
		_, err := groupStmt.Exec(groupUUID, groupName)
		if err != nil {
			slog.Error(err.Error())
		}
		groupUUID++
	}
}

type member struct {
	uuid         string
	receiveEmail bool
	archived     bool
	familyUUID   string
	familyGroup  string
	groupUUIDs   []string
}

var members = map[string]member{
	"11": {uuid: "11", receiveEmail: true, archived: false, familyUUID: "101", familyGroup: "fam1", groupUUIDs: []string{"1001", "1002"}},
	"12": {uuid: "12", receiveEmail: true, archived: false, familyUUID: "101", familyGroup: "fam1", groupUUIDs: []string{"1002", "1003"}},
	"13": {uuid: "13", receiveEmail: false, archived: false, familyUUID: "102", familyGroup: "fam2", groupUUIDs: []string{"1001", "1002"}},
	"14": {uuid: "14", receiveEmail: true, archived: false, familyUUID: "102", familyGroup: "fam2", groupUUIDs: []string{"1002", "1003"}},
	"15": {uuid: "15", receiveEmail: false, archived: true, familyUUID: "102", familyGroup: "fam2", groupUUIDs: []string{"1002", "1004"}},
}

func (r *Repo) InsertMembers() {
	nofMembers := 72
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
		familyGroup := ""
		groupUUIDs := []string{}
		if m, ok := members[strconv.Itoa(memberUUID)]; ok {
			receiveEmail = m.receiveEmail
			archived = m.archived
			familyUUID = m.familyUUID
			familyGroup = m.familyGroup
			groupUUIDs = m.groupUUIDs
		}

		member := entity.NewMember(strconv.Itoa(memberUUID), memberId, name, dob, personnummer, email, mobile, address, "", "", util.String2Time("2020-01-01"), time.Time{}, receiveEmail, false, false, archived, status, familyUUID, familyGroup)
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
