package repo

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/util"
)

const dbName = "sqlite3"

type Repo struct {
	db *sql.DB
}

func NewRepo() *Repo {
	var dbFile string

	if os.Getenv("APP_ENV") == "production" {
		dbFile = filepath.Join(os.Getenv(constants.PROD_DB_DIR_KEY), os.Getenv(constants.DB_NAME_KEY))
	} else {
		dbFile = filepath.Join(os.Getenv(constants.DEV_DB_DIR_KEY), os.Getenv(constants.DB_NAME_KEY))
	}

	fmt.Println("[REPO]:NewRepo", "dbFile:", dbFile)

	db, err := sql.Open(dbName, dbFile)
	if err != nil {
		slog.Error(err.Error())
	}
	// 	defer db.Close()
	return &Repo{db}
}

func (r *Repo) Close() {
	r.db.Close()
}

func (r *Repo) GetDB() *sql.DB {
	return r.db
}

func (r *Repo) IsAuthenticated(username string, password string) bool {
	var name string
	var pw string

	err := r.db.QueryRow("SELECT name, password FROM users WHERE name = ?", username).Scan(&name, &pw)
	if err != nil {
		slog.Error(err.Error(), "name", username)
		return false
	}
	fmt.Println("[REPO]:IsAuthenticated", "name:", name, "pw:", pw)

	return util.ValidatePassword(password, pw)
}

func (r *Repo) InitRoot() {
	stmt, err := r.db.Prepare("INSERT INTO users(name,password) values(?, ?)")
	if err != nil {
		slog.Error(err.Error())
	}
	hpw, _ := util.HashPassword("joe")
	_, err = stmt.Exec("root", hpw)
	if err != nil {
		slog.Error(err.Error())
	}
}

func (r *Repo) InitDDL() {
	var sqlStmts = make(map[string]string)

	sqlStmts["drop_user"] = `DROP TABLE IF EXISTS users;`
	sqlStmts["user"] = `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	for name, sqlStmt := range sqlStmts {
		_, err := r.db.Exec(sqlStmt)
		if err != nil {
			slog.Error(fmt.Sprintf("Error in stmt %q:%s\n%s\n", err, name, sqlStmt))
		}
	}
}
