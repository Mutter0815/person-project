package db

import (
	"log"
	"person-project/models"
)

func Migrate() {
	err := DB.AutoMigrate(&models.Person{})
	if err != nil {
		log.Fatal("Ошибка миграции", err)
	}
	log.Println("Миграции успешно применены")
}
