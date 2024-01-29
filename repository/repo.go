package repo

import (
	"database/sql"
	"log/slog"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/michelemendel/dmtmms/constants"
)

const dbName = "sqlite3"

type Repo struct {
	DB *sql.DB
}

func NewRepo() *Repo {
	var dbFile string

	if os.Getenv("APP_ENV") == "production" {
		dbFile = filepath.Join(os.Getenv(constants.PROD_DB_DIR_KEY), os.Getenv(constants.DB_NAME_KEY))
	} else {
		dbFile = filepath.Join(os.Getenv(constants.DEV_DB_DIR_KEY), os.Getenv(constants.DB_NAME_KEY))
	}

	db, err := sql.Open(dbName, dbFile)
	if err != nil {
		slog.Error(err.Error())
	}
	return &Repo{db}
}

func (r *Repo) Close() {
	r.DB.Close()
}

func (r *Repo) GetDB() *sql.DB {
	return r.DB
}
