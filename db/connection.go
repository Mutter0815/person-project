package db

import (
	"fmt"
	"log"
	"person-project/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {

	cfg := config.Cfg
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s  sslmode=disable", config.Cfg.DB.Host, config.Cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name)
	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных %v", err)
	}

	log.Println("Успешное подключение к базе данных")
}
