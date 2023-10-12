package globals

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var Secret = []byte("secret")

const Userkey = "user"

var (
	USERNAME string
	PASSWORD string
	NETWORK  string
	HOST     string
	PORT     int
	DATABASE string
)

func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	USERNAME = os.Getenv("_USERNAME")
	PASSWORD = os.Getenv("PASSWORD")
	NETWORK = os.Getenv("NETWORK")
	HOST = os.Getenv("HOST")
	_PORT := os.Getenv("PORT")
	DATABASE = os.Getenv("DATABASE")
	PORT, _ = strconv.Atoi(_PORT)
}
