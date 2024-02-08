package main

import (
	"flag"
	"fmt"

	repo "github.com/michelemendel/dmtmms/repository"
)

func main() {
	db := flag.String("db", "", "reset")
	drop := flag.String("drop", "", "tables")
	create := flag.String("create", "", "tables")
	insert := flag.String("insert", "", "users, members_groups")
	show := flag.String("show", "", "users")

	flag.Parse()

	repo := repo.NewRepo()
	defer repo.Close()
	// repo.DBConfig()

	if *db != "" {
		switch *db {
		case "reset":
			repo.DropTables()
			repo.CreateTables()
			repo.CreateIndexes()
			repo.InsertUsers()
			repo.InsertMembersGroups()
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
			repo.InsertMembersGroups()
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
