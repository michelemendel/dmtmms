package repo

import (
	"fmt"
	"log/slog"

	"github.com/michelemendel/dmtmms/util"
)

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

func (r *Repo) RunDDL() {
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
