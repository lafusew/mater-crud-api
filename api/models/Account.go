package models

import "github.com/jinzhu/gorm"

type Account struct {
	gorm.Model
	Balance uint `json:"balance"`
}

type TeamAccount struct {
	Account
	TeamID uint `gorm:"not null" json:"team_id"`
}

type UserAccount struct {
	Account
	UserID uint `gorm:"not null" json:"user_id"`
}
