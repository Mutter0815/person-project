package models

import (
	"time"
)

type Gender string

const (
	GenderMale   string = "male"
	GenderFemale string = "female"
)

type Person struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created_at" gorm:"index"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name" binding:"required" gorm:"index"`
	Surname     string    `json:"surname" binding:"required" gorm:"index"`
	Patronymic  *string   `json:"patronymic,omitempty"`
	Age         int       `json:"age" gorm:"index"`
	Gender      Gender    `json:"gender" validate:"oneof=male female" gorm:"index"`
	Nationality string    `json:"nationality"`
}
