package util

import (
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var once sync.Once

func InitEnv() {
	once.Do(InitEnvExec)
}

func InitEnvExec() {
	envFile := ".env"
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatal("[utils]:error loading env file:", envFile)
	}
}

func GenerateUuid() string {
	return uuid.NewString()
}
