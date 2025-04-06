package models

import "gorm.io/gorm"

type Person struct {
	gorm.Model
	Name        string  `json:"name" binding:"required"`
	Surname     string  `json:"surname" binding:"required"`
	Patronymic  *string `json:"patronymic"`
	Age         int     `json:"age"`
	Gender      string  `json:"gender"`
	Nationality string  `json:"nationality"`
}
