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

	sqlStmts["drop_user"] = `DROP TABLE IF EXISTS users;`
	sqlStmts["drop_members"] = `DROP TABLE IF EXISTS members;`
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
		id TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		dob REAL,
		email TEXT,
		mobile TEXT,
		status TEXT NOT NULL,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER
	); `

	sqlStmts["create_groups"] = `
	CREATE TABLE IF NOT EXISTS groups (
		uuid TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER
	);`

	// many-to-many between members and groups
	sqlStmts["create_members_groups"] = `
	CREATE TABLE IF NOT EXISTS members_groups (
		member_uuid TEXT NOT NULL,
		group_uuid TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		updated_at INTEGER,
		primary key (member_uuid, group_uuid)
		FOREIGN KEY(member_uuid) REFERENCES members(uuid),
		FOREIGN KEY(group_uuid) REFERENCES groups(uuid)
	); `

	r.runStatements(sqlStmts)
}
func (r *Repo) CreateIndexes() {
	var sqlStmts = make(map[string]string)

	sqlStmts["create_index_members_uuid"] = `
	CREATE UNIQUE INDEX IF NOT EXISTS idx_members_uuid ON members(uuid);`

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

func (r *Repo) InsertMembersGroups() {
	memberStmt, _ := r.DB.Prepare("INSERT INTO members(uuid, id, name, dob, email, mobile, status) values(?, ?, ?, julianday(?), ?, ?, ?)")
	groupStmt, _ := r.DB.Prepare("INSERT INTO groups(uuid, name, type) values(?, ?, ?)")
	memberGroupStmt, _ := r.DB.Prepare("INSERT INTO members_groups(member_uuid, group_uuid, role) values(?, ?, ?)")

	type member struct {
		name   string
		email  string
		mobile string
	}

	members := []member{
		{"mem1", "mem1@t.c", "12377891"},
		{"mem2", "mem2@t.c", "12377892"},
		{"mem3", "mem3@t.c", "12377893"},
		{"mem4", "mem4@t.c", "12377894"},
	}

	userId := 1
	memberUUID := 10
	var status entity.MemberStatus = entity.MemberStatusActive
	dob := util.String2Time("1965-07-22")
	for _, member := range members {
		_, err := memberStmt.Exec(strconv.Itoa(memberUUID), strconv.Itoa(userId), member.name, dob, member.email, member.mobile, status)
		if err != nil {
			slog.Error(err.Error())
		}
		userId++
		memberUUID++
		dob = dob.AddDate(0, 1, 1)
	}

	groupUUID := 100
	type group struct {
		name string
		typ  string
	}
	groups := []group{
		{"fam1", "fam"},
		{"fam2", "fam"},
		{"org1", "org"},
		{"org2", "org"},
	}
	for _, group := range groups {
		_, err := groupStmt.Exec(groupUUID, group.name, group.typ)
		if err != nil {
			slog.Error(err.Error())
		}
		groupUUID++
	}

	// mem1,mem2 is in fam1
	_, _ = memberGroupStmt.Exec(10, 100, "parent")
	_, _ = memberGroupStmt.Exec(11, 100, "child")
	// mem3,mem4 is in fam2
	_, _ = memberGroupStmt.Exec(12, 101, "parent")
	_, _ = memberGroupStmt.Exec(13, 101, "child")
	// mem1,mem3 is in org1 and org2
	_, _ = memberGroupStmt.Exec(10, 102, "leader")
	_, _ = memberGroupStmt.Exec(10, 103, "finance")
	_, _ = memberGroupStmt.Exec(12, 102, "house")
	_, _ = memberGroupStmt.Exec(12, 103, "children's activities")
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
