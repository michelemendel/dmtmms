package main

import (
	"fmt"
	"os"

	consts "github.com/michelemendel/dmtmms/constants"
	repo "github.com/michelemendel/dmtmms/repository"
)

func main() {
	fmt.Println("IMPORT DMT DATA")

	env := os.Getenv(consts.ENV_APP_ENV_KEY)
	webServerPort := os.Getenv(consts.ENV_WEB_SERVER_PORT_KEY)

	fmt.Printf("ENVIRONMENT:\nmode:%s\nwebServerPort:%s\n", env, webServerPort)

	r := repo.NewRepo()
	defer r.DB.Close()

	r.ImportData()

}
