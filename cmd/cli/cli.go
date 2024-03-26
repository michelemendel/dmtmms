package main

import (
	"flag"
	"fmt"

	repo "github.com/michelemendel/dmtmms/repository"
)

func main() {
	db := flag.String("db", "", "reset, reset_no_auto_id")
	drop := flag.String("drop", "", "tables")
	create := flag.String("create", "", "tables, tables_no_auto_id")
	insert := flag.String("insert", "", "users, members_groups")
	show := flag.String("show", "", "users")

	flag.Parse()

	repo := repo.NewRepo()
	defer repo.Close()
	repo.DB.Exec("PRAGMA foreign_keys = OFF")

	if *db != "" {
		switch *db {
		case "reset":
			repo.DropTables()
			repo.CreateTables()
			repo.CreateTriggers(true)
			repo.CreateIndexes()
			repo.InsertUsers()
		// repo.InsertFamilies()
		// repo.InsertGroups()
		// repo.InsertMembers()
		case "reset_no_auto_id":
			repo.DropTables()
			repo.CreateTables()
			repo.CreateTriggers(false)
			repo.CreateIndexes()
			repo.InsertUsers()
		default:
			fmt.Println("no op specified")
		}
	}

	if *drop != "" {
		switch *drop {
		case "tables":
			repo.DropTables()
		default:
			fmt.Println("no op specified")
		}
	}

	if *create != "" {
		switch *create {
		case "tables":
			repo.CreateTables()
			repo.CreateTriggers(true)
			repo.CreateIndexes()
		case "tables_no_auto_id":
			repo.CreateTables()
			repo.CreateTriggers(false)
			repo.CreateIndexes()
		default:
			fmt.Println("no op specified")
		}
	}

	if *insert != "" {
		switch *insert {
		case "users":
			repo.InsertUsers()
		case "members_groups":
			repo.InsertFamilies()
			repo.InsertGroups()
			repo.InsertMembers()
		default:
			fmt.Println("no op specified")
		}
	}

	if *show != "" {
		switch *show {
		case "users":
			repo.ShowUsers()
		default:
			fmt.Println("no op specified")
		}
	}
}
