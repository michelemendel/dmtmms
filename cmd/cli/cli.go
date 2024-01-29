package main

import (
	"flag"
	"fmt"

	repo "github.com/michelemendel/dmtmms/repository"
)

func main() {
	op := flag.String("op", "", "migrate,initdata,users")
	flag.Parse()
	repo := repo.NewRepo()
	defer repo.Close()

	fmt.Println("op:", *op)

	switch *op {
	case "migrate":
		repo.DBConfig()
		repo.RunDDL()
	case "initdata":
		repo.InitRootUser()
	case "users":
		repo.GetUsers()
	default:
		fmt.Println("no op specified")
	}
}
