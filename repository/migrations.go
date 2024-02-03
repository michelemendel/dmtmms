package repo

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/michelemendel/dmtmms/util"
)

func (r *Repo) DBConfig() {
	r.DB.Exec("PRAGMA journal_mode = WAL")
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

	sqlStmts["create_members"] = `
	CREATE TABLE IF NOT EXISTS members (
		uuid TEXT PRIMARY KEY,
		id TEXT NOT NULL UNIQUE,
		name TEXT NOT NULL,
		date_of_birth REAL
		created_at INTEGER,
		updated_at INTEGER
	); `

	sqlStmts["create_groups"] = `
	CREATE TABLE IF NOT EXISTS groups (
		uuid TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		type TEXT NOT NULL,
		created_at INTEGER,
		updated_at INTEGER
	);`

	// many-to-many between members and groups
	sqlStmts["create_members_groups"] = `
	CREATE TABLE IF NOT EXISTS members_groups (
		member_uuid TEXT NOT NULL,
		group_uuid TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at INTEGER,
		updated_at INTEGER,
		primary key (member_uuid, group_uuid)
	); `

	r.runStatements(sqlStmts)
}

func (r *Repo) RunDML() {}

func (r *Repo) InsertUsers() {
	// stmt, err := r.DB.Prepare("INSERT INTO users(name,password,role,date_of_birth) values(?, ?, ?, julianday(?))")
	stmt, err := r.DB.Prepare("INSERT INTO users(name,password,role) values(?, ?, ?)")
	if err != nil {
		slog.Error(err.Error())
	}

	for _, user := range []string{"root", "abe", "bob"} {
		hpw, _ := util.HashPassword(user)
		_, err = stmt.Exec(user, hpw, "admin")
		if err != nil {
			slog.Error(err.Error())
		}
	}
}

func (r *Repo) InsertMembersGroups() {
	memberStmt, _ := r.DB.Prepare("INSERT INTO members(uuid,id,name,date_of_birth) values(?, ?, ?, julianday(?))")
	groupStmt, _ := r.DB.Prepare("INSERT INTO groups(uuid,name,type) values(?, ?, ?)")
	// memberGroupStmt, err := r.DB.Prepare("INSERT INTO groups(member_uuid,group_uuid,role,type) values(?, ?, ?, ?)")

	userId := 1
	for _, member := range []string{"mem1", "mem2", "mem3"} {
		memberUUID := util.GenerateUUID()
		// groupUUID := util.GenerateUUID()
		fmt.Println("member: ", member)

		_, err := memberStmt.Exec(memberUUID, strconv.Itoa(userId), member, "1965-07-22")
		if err != nil {
			slog.Error(err.Error())
		}
		userId++
	}

	for _, group := range []string{"fam1", "fam2"} {
		groupUUID := util.GenerateUUID()
		// fmt.Println("member: ", member)

		_, err := groupStmt.Exec(groupUUID, group, "family")
		if err != nil {
			slog.Error(err.Error())
		}
		userId++
	}
}

func (r *Repo) ShowUsers() error {
	var name string
	var password string
	var role string
	var createdAt string
	// -- AND date_of_birth < julianday('1965-07-22')
	// rows, err := r.DB.Query(`SELECT name,date(date_of_birth),datetime(created_at,'LOCALTIME') FROM users;`)
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
