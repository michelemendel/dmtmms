package repo

import (
	"fmt"
	"log/slog"

	"github.com/michelemendel/dmtmms/util"
)

func (r *Repo) DBConfig() {
	r.DB.Exec("PRAGMA journal_mode = WAL")
	// This doesn't work
	r.DB.Exec("PRAGMA busy_timeout = 5000")
}

func (r *Repo) RunDDL() {
	var sqlStmts = make(map[string]string)

	// Users using the application
	// TODO: remove date_of_birth
	sqlStmts["drop_user"] = `DROP TABLE IF EXISTS users;`
	sqlStmts["user"] = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL,
		created_at INTEGER DEFAULT CURRENT_TIMESTAMP,
		date_of_birth REAL
	);
	`

	for name, sqlStmt := range sqlStmts {
		_, err := r.DB.Exec(sqlStmt)
		if err != nil {
			slog.Error(fmt.Sprintf("Error in stmt %q:%s\n%s\n", err, name, sqlStmt))
		}
	}
}

func (r *Repo) RunDML() {}

func (r *Repo) InitRootUser() {
	stmt, err := r.DB.Prepare("INSERT INTO users(name,password,role,date_of_birth) values(?, ?, ?, julianday(?))")
	if err != nil {
		slog.Error(err.Error())
	}
	hpw, _ := util.HashPassword("joe")
	_, err = stmt.Exec("root", hpw, "admin", "1965-07-24")
	if err != nil {
		slog.Error(err.Error())
	}
}

func (r *Repo) GetUsers() {
	var name string
	var dateOfBirth string
	var createdAt string
	// -- AND date_of_birth < julianday('1965-07-22')
	rows, err := r.DB.Query(`SELECT name,date(date_of_birth),datetime(created_at,'LOCALTIME') FROM users;`)
	if err != nil {
		slog.Error(err.Error())
	}
	defer rows.Close()
	tabs := "%s\t%s\t%s\n"
	fmt.Printf(tabs, "name", "dateOfBirth", "createdAt")
	for rows.Next() {
		err := rows.Scan(&name, &dateOfBirth, &createdAt)
		if err != nil {
			slog.Error(err.Error())
		}
		fmt.Printf(tabs, name, dateOfBirth, createdAt)
	}

	// fmt.Println("[REPO]:GetUser:", name, dateOfBirth, createdAt)
}
