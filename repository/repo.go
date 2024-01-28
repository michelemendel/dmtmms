package repo

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/michelemendel/dmtmms/constants"
	"github.com/michelemendel/dmtmms/entity"
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

func (r *Repo) GetUser(username string) (entity.User, error) {
	var name string
	var pw string

	err := r.db.QueryRow("SELECT name, password FROM users WHERE name = ?", username).Scan(&name, &pw)
	if err != nil {
		slog.Error(err.Error(), "name", username)
		return entity.User{}, err
	}
	fmt.Println("[REPO]:IsAuthenticated", "name:", name, "pw:", pw)
	return entity.NewUser(name, pw), nil
}
