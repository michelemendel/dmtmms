package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/michelemendel/dmtmms/constants"
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

func Date2String(t time.Time) string {
	return t.Format(constants.DATE_FRMT)
}

func String2Date(s string) time.Time {
	t, err := time.Parse(constants.DATE_FRMT, s)
	if err != nil {
		// slog.Error("Error parsing date", "error", err)
		return time.Time{}
	}
	return t
}

func CalculateAge(DOB time.Time) int {
	today := time.Now()
	age := today.Year() - DOB.Year()
	if today.YearDay() < DOB.YearDay() {
		age--
	}
	return age
}

func DateTime2String(t time.Time) string {
	loc, err := time.LoadLocation("Europe/Oslo")
	if err != nil {
		slog.Error("Error loading location", "error", err)
		return ""
	}
	t = t.In(loc)
	return t.Format(constants.DATE_TIME_FRMT)
}

func String2DateTime(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	t, err := time.Parse(constants.DATE_TIME_FRMT, s)
	if err != nil {
		slog.Error("Error parsing date time", "error", err)
		return time.Time{}
	}
	return t
}

func GetYearFromAge(age string) string {
	year := time.Now().Year() - String2Int(age)
	return Int2String(year)
}

func String2Int(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func Int2String(i int) string {
	return strconv.Itoa(i)
}

func String2Bool(s string) bool {
	return s != ""
}

func Bool2String(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// Check if it is a hxr request
func IsHXR(c echo.Context) bool {
	if c.QueryParam("l") == "ok" {
		return false
	}

	return c.Request().Header.Get("HX-Request") == "true"
}

func PP(s any) {
	res, err := PrettyStruct(s)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(res)
}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}

type Age struct {
	Val  string
	Name string
}

func MakeAges() []Age {
	ages := []Age{{Val: "", Name: ""}}

	for i := 0; i <= 5; i++ {
		ages = append(ages, Age{Val: fmt.Sprintf("%d", i), Name: fmt.Sprintf("%d", i)})
	}

	for i := 6; i <= 18; i++ {
		ages = append(ages, Age{Val: fmt.Sprintf("%d", i), Name: fmt.Sprintf("%d, %d kl", i, i-5)})
	}

	for i := 19; i <= 120; i++ {
		ages = append(ages, Age{Val: fmt.Sprintf("%d", i), Name: fmt.Sprintf("%d", i)})
	}

	return ages
}
