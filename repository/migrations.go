package repo

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/michelemendel/dmtmms/entity"
	"github.com/michelemendel/dmtmms/util"
)

func (r *Repo) DBConfig() {
	r.DB.Exec("PRAGMA journal_mode = WAL")
	r.DB.Exec("PRAGMA foreign_keys = ON")
	// This doesn't work
	r.DB.Exec("PRAGMA busy_timeout = 5000")
}

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
	// state = [active, archived, tobedeleted]
	sqlStmts["create_members"] = `
	CREATE TABLE IF NOT EXISTS members (
		uuid TEXT PRIMARY KEY,
		family_uuid TEXT,
		id TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		dob REAL,
		email TEXT,
		mobile TEXT,
		status TEXT NOT NULL,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER,
		FOREIGN KEY(family_uuid) REFERENCES families(uuid)
	); `

	sqlStmts["create_families"] = `
	CREATE TABLE IF NOT EXISTS families (
		uuid TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER
	);`

	sqlStmts["create_groups"] = `
	CREATE TABLE IF NOT EXISTS groups (
		uuid TEXT PRIMARY KEY,
		name TEXT NOT NULL,
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

func (r *Repo) InsertMembersGroups() {
	memberStmt, _ := r.DB.Prepare("INSERT INTO members(uuid, id, name, dob, email, mobile, status) values(?, ?, ?, julianday(?), ?, ?, ?)")
	familyStmt, _ := r.DB.Prepare("INSERT INTO families(uuid, name) values(?, ?)")
	groupStmt, _ := r.DB.Prepare("INSERT INTO groups(uuid, name) values(?, ?)")

	memberIdPrefix := "99"
	namePrefix := "mem"
	phonePrefix := "12377"
	nrStart := 100
	memberUUID := 11
	nofMembers := 50
	var status entity.MemberStatus = entity.MemberStatusActive
	dob := util.String2Time("1980-02-01")
	for i := 0; i < nofMembers; i++ {
		memberId := memberIdPrefix + strconv.Itoa(nrStart+i)
		name := namePrefix + strconv.Itoa(memberUUID)
		email := name + "@test.com"
		mobile := phonePrefix + strconv.Itoa(nrStart+i)
		_, err := memberStmt.Exec(strconv.Itoa(memberUUID), memberId, name, dob, email, mobile, status)
		if err != nil {
			slog.Error(err.Error())
		}
		memberUUID++
		dob = dob.AddDate(0, 1, 1)
	}

	// Families

	familyUUID := 101
	families := []string{"fam1", "fam2"}
	for _, familyName := range families {
		_, err := familyStmt.Exec(familyUUID, familyName)
		if err != nil {
			slog.Error(err.Error())
		}
		familyUUID++
	}

	// Groups

	groupUUID := 1001
	groups := []string{"org1", "org2", "org3"}
	for _, groupName := range groups {
		_, err := groupStmt.Exec(groupUUID, groupName)
		if err != nil {
			slog.Error(err.Error())
		}
		groupUUID++
	}

	// Family relationships

	// Update members family relation
	famUpdStmt, err := r.DB.Prepare("UPDATE members SET family_uuid=? WHERE uuid=?")

	if err != nil {
		fmt.Println("error in prepared stmt for update", err.Error())
	}
	// mem1,mem2 is in fam1
	famUpdStmt.Exec(101, 11)
	famUpdStmt.Exec(101, 12)
	// mem3,mem4 is in fam2
	famUpdStmt.Exec(102, 13)
	famUpdStmt.Exec(102, 14)

	// Group relationships
	memberGroupStmt, _ := r.DB.Prepare("INSERT INTO members_groups(member_uuid, group_uuid) values(?, ?)")

	// mem1,mem3 is in org1 and org2
	_, _ = memberGroupStmt.Exec(11, 1001)
	_, _ = memberGroupStmt.Exec(11, 1002)
	_, _ = memberGroupStmt.Exec(13, 1001)
	_, _ = memberGroupStmt.Exec(13, 1002)
	// mem4 is in org2 and org3
	_, _ = memberGroupStmt.Exec(14, 1002)
	_, _ = memberGroupStmt.Exec(14, 1003)
	// mem5 is in org2 and org3
	_, _ = memberGroupStmt.Exec(15, 1002)
	_, _ = memberGroupStmt.Exec(15, 1003)
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
