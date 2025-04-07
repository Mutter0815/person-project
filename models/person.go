package models

import (
	"time"
)

type Person struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name" binding:"required"`
	Surname     string    `json:"surname" binding:"required"`
	Patronymic  *string   `json:"patronymic,omitempty"`
	Age         int       `json:"age"`
	Gender      string    `json:"gender"`
	Nationality string    `json:"nationality"`
}

type CreatePersonRequest struct {
	Name       string  `json:"name" example:"Ivan"`
	Surname    string  `json:"surname" example:"Ivanov"`
	Patronymic *string `json:"patronymic" example:"Ivanovich"`
}
