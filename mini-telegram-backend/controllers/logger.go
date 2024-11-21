package controllers

import (
	"mini-telegram/config"
	"mini-telegram/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogHandler(c *gin.Context) {
	var logEntry models.LogEntry

	if err := c.ShouldBindJSON(&logEntry); err != nil {
		config.Logger.Warnf("Invalid log entry data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	go func() {
		newLog := models.LogEntry{
			Message:  logEntry.Message,
			Stack:    logEntry.Stack,
			UserID:   logEntry.UserID,
			Username: logEntry.Username,
		}

		if err := config.DB.Create(&newLog).Error; err != nil {
			config.Logger.Errorf("Failed to save log entry: %v", err)
		} else {
			config.Logger.Infof("Log entry saved: %v", logEntry.Message)
		}
	}()

	config.Logger.Infof("Log entry received: %s", logEntry.Message)
	c.JSON(http.StatusOK, gin.H{"message": "Error logged successfully"})
}
