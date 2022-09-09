package module

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
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

func LoadEnv(debug bool) Environments {
	var err error
	// LOADS .env file
	if debug {
		err = godotenv.Load(".env")
	} else {
		err = godotenv.Load(".env_prod")
	}

	if err != nil {
		log.Fatal("Err: Loading .env failed.")
	}

	env := new(Environments)

	// Load environment values
	env.AppPort = os.Getenv("APP_PORT")
	env.AppHost = os.Getenv("APP_HOST")
	env.YEAR, _ = strconv.Atoi(os.Getenv("YEAR"))
	env.LineChannelAccessToken = os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	env.LineChannelSecret = os.Getenv("LINE_CHANNEL_SECRET")
	env.LineAdminUid = os.Getenv("LINE_ADMIN_UID")
	env.LineMockUid = os.Getenv("LINE_MOCK_UID")
	env.DbHost = os.Getenv("DB_HOST")
	env.DbPort = os.Getenv("DB_PORT")
	env.DbUser = os.Getenv("DB_USER")
	env.DbPass = os.Getenv("DB_PASS")
	env.DbName = os.Getenv("DB_NAME")
	env.KuwikiEndpoint = os.Getenv("KUWIKI_ENDPOINT")
	env.KuwikiAccessToken = os.Getenv("KUWIKI_ACCESS_TOKEN")
	env.GmailId = os.Getenv("GMAIL_ID")
	env.GmailPassword = os.Getenv("GMAIL_PASSWORD")

	return *env
}
