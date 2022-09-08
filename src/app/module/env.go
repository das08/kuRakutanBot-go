package module

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Environments struct {
	APP_PORT                  string
	APP_HOST                  string
	YEAR                      int
	LINE_CHANNEL_ACCESS_TOKEN string
	LINE_CHANNEL_SECRET       string
	LINE_ADMIN_UID            string
	LINE_MOCK_UID             string
	DB_HOST                   string
	DB_PORT                   string
	DB_USER                   string
	DB_PASS                   string
	DB_NAME                   string
	DB_COLLECTION             Collections
	KUWIKI_ENDPOINT           string
	KUWIKI_ACCESS_TOKEN       string
	GMAIL_ID                  string
	GMAIL_PASSWORD            string
}

type Collection = string

type Collections struct {
	User         Collection
	Rakutan      Collection
	Favorites    Collection
	Verification Collection
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
	env.APP_PORT = os.Getenv("APP_PORT")
	env.APP_HOST = os.Getenv("APP_HOST")
	env.YEAR, _ = strconv.Atoi(os.Getenv("YEAR"))
	env.LINE_CHANNEL_ACCESS_TOKEN = os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	env.LINE_CHANNEL_SECRET = os.Getenv("LINE_CHANNEL_SECRET")
	env.LINE_ADMIN_UID = os.Getenv("LINE_ADMIN_UID")
	env.LINE_MOCK_UID = os.Getenv("LINE_MOCK_UID")
	env.DB_HOST = os.Getenv("DB_HOST")
	env.DB_PORT = os.Getenv("DB_PORT")
	env.DB_USER = os.Getenv("DB_USER")
	env.DB_PASS = os.Getenv("DB_PASS")
	env.DB_NAME = os.Getenv("DB_NAME")
	env.DB_COLLECTION.User = os.Getenv("DB_COL_USER")
	env.DB_COLLECTION.Rakutan = os.Getenv("DB_COL_RAKUTAN")
	env.DB_COLLECTION.Favorites = os.Getenv("DB_COL_FAV")
	env.DB_COLLECTION.Verification = os.Getenv("DB_COL_VER")
	env.KUWIKI_ENDPOINT = os.Getenv("KUWIKI_ENDPOINT")
	env.KUWIKI_ACCESS_TOKEN = os.Getenv("KUWIKI_ACCESS_TOKEN")
	env.GMAIL_ID = os.Getenv("GMAIL_ID")
	env.GMAIL_PASSWORD = os.Getenv("GMAIL_PASSWORD")

	return *env
}
