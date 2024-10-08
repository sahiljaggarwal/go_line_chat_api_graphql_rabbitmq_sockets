package db

import (
	"fmt"
	"line/src/configs/env"
	"line/src/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		env.DB_Host,
		env.DB_User,
		env.DB_Password,
		env.DB_Name,
		env.DB_Port,
	)
	// log.Print(dsn)
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error on connecting database: ", err)
		return
	}
	err = DB.AutoMigrate(&models.User{}, &models.Conversation{}, &models.Message{})
	if err != nil {
		log.Fatal("Failed to auto migrate the user model: ", err)
	}
	fmt.Println("Database connected successfully.")

}
