package module

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Environments struct {
	LINE_CHANNEL_ACCESS_TOKEN string
	LINE_CHANNEL_SECRET       string
	LINE_ADMIN_UID            string
	DB_HOST                   string
	DB_PORT                   string
	DB_USER                   string
	DB_PASS                   string
	DB_NAME                   string
	DB_COLLECTION             string
	KUWIKI_ENDPOINT           string
	KUWIKI_ACCESS_TOKEN       string
}

func LoadEnv(debug bool) Environments {
	var err error
	// LOADS .env file
	if debug {
		err = godotenv.Load("../../.env")
	} else {
		err = godotenv.Load("../../.env_prod")
	}

	if err != nil {
		log.Fatal("Err: Loading .env failed.")
	}

	env := new(Environments)

	// Load environment values
	env.LINE_CHANNEL_ACCESS_TOKEN = os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	env.LINE_CHANNEL_SECRET = os.Getenv("LINE_CHANNEL_SECRET")
	env.LINE_ADMIN_UID = os.Getenv("LINE_ADMIN_UID")
	env.DB_HOST = os.Getenv("DB_HOST")
	env.DB_PORT = os.Getenv("DB_PORT")
	env.DB_USER = os.Getenv("DB_USER")
	env.DB_PASS = os.Getenv("DB_PASS")
	env.DB_NAME = os.Getenv("DB_NAME")
	env.DB_COLLECTION = os.Getenv("DB_COLLECTION")
	env.KUWIKI_ENDPOINT = os.Getenv("KUWIKI_ENDPOINT")
	env.KUWIKI_ACCESS_TOKEN = os.Getenv("KUWIKI_ACCESS_TOKEN")

	return *env
}
