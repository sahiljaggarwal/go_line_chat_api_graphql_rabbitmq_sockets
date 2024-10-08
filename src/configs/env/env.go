package env

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	PORT string
	HOST string

	DB_Host     string
	DB_User     string
	DB_Password string
	DB_Name     string
	DB_Port     string

	SECRET_KEY   string
	RABBITMQ_URL string
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error on loading .env file")
	}

	PORT = os.Getenv("PORT")
	HOST = os.Getenv("HOST")
	DB_Host = os.Getenv("DB_HOST")
	DB_User = os.Getenv("DB_USER")
	DB_Password = os.Getenv("DB_PASSWORD")
	DB_Name = os.Getenv("DB_NAME")
	DB_Port = os.Getenv("DB_PORT")

	SECRET_KEY = os.Getenv("SECRET_KEY")

	RABBITMQ_URL = os.Getenv("RABBITMQ_URL")
}
