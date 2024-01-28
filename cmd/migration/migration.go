package main

import (
	"flag"
	"fmt"

	repo "github.com/michelemendel/dmtmms/repository"
)

func main() {
	op := flag.String("op", "", "migrate,fill")
	flag.Parse()
	repo := repo.NewRepo()
	defer repo.Close()

	fmt.Println("[MIGRATION]", "op:", *op)

	switch *op {
	case "migrate":
		repo.RunDDL()
	case "fill":
		repo.InitRoot()
	default:
		fmt.Println("no op specified")
	}
}
