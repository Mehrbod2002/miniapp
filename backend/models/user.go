package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	TelegramID int64  `json:"telegram_id" gorm:"unique"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	Points     int    `json:"points"`
}
