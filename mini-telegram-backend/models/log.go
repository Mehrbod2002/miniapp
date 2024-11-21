package models

import "gorm.io/gorm"

type LogEntry struct {
	gorm.Model
	Message  string `json:"message"`
	Stack    string `json:"stack"`
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
}
