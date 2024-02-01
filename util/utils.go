package util

import (
	"encoding/base64"
	"log"
	"log/slog"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
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

func GenerateUUID() string {
	return uuid.NewString()
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		slog.Error("Error hashing password", "error", err)
		// Just create some random data, so the password is not empty and can be compromised.
		return string(GenerateUUID()), err
	}
	return string(bytes), nil
}

func ValidatePassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Generate a random password with n charactes.
func GeneratePassword() string {
	aLongString := base64.StdEncoding.EncodeToString([]byte(GenerateUUID()))
	return strings.ToLower(aLongString[0:8])
}
