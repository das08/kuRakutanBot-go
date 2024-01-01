package module

import (
	"log"
	"os"
	"strconv"
)

type Environments struct {
	AppPort                string
	AppHost                string
	YEAR                   int
	LineChannelAccessToken string
	LineChannelSecret      string
	LineAdminUid           string
	LineMockUid            string
	DbHost                 string
	DbPort                 string
	DbUser                 string
	DbPass                 string
	DbName                 string
	KuwikiEndpoint         string
	KuwikiAccessToken      string
	GmailId                string
	GmailPassword          string
}

func getEnv(key string) string {
	// Get environment value
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Error: %s is not set.", key)
	}
	log.Println(key, "=", value)
	return value
}

func LoadEnv() Environments {
	env := new(Environments)

	// Load environment values
	env.AppPort = getEnv("APP_PORT")
	env.AppHost = getEnv("APP_HOST")
	env.YEAR, _ = strconv.Atoi(getEnv("YEAR"))
	env.LineChannelAccessToken = getEnv("LINE_CHANNEL_ACCESS_TOKEN")
	env.LineChannelSecret = getEnv("LINE_CHANNEL_SECRET")
	env.LineAdminUid = getEnv("LINE_ADMIN_UID")
	env.LineMockUid = getEnv("LINE_MOCK_UID")
	env.DbHost = getEnv("DB_HOST")
	env.DbPort = getEnv("DB_PORT")
	env.DbUser = getEnv("DB_USER")
	env.DbPass = getEnv("DB_PASS")
	env.DbName = getEnv("DB_NAME")
	env.KuwikiEndpoint = getEnv("KUWIKI_ENDPOINT")
	env.KuwikiAccessToken = getEnv("KUWIKI_ACCESS_TOKEN")
	env.GmailId = getEnv("GMAIL_ID")
	env.GmailPassword = getEnv("GMAIL_PASSWORD")

	return *env
}
