package controllers

import (
	"mini-telegram/config"
	"mini-telegram/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Register(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		config.Logger.Warnf("Invalid input data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}

	var user models.User
	if err := config.DB.Where("telegram_id = ?", input.TelegramID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User doesn't exist, create a new one
			input.Points = 0
			if err := config.DB.Create(&input).Error; err != nil {
				config.Logger.Errorf("Database error while registering user: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not register user"})
				return
			}
			user = input
		} else {
			config.Logger.Errorf("Database error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	config.Logger.Infof("User registered/login successful: Telegram ID %d", user.TelegramID)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"points":  user.Points,
	})
}

func GetPoints(c *gin.Context) {
	telegramID := c.Param("telegram_id")
	config.Logger.Infof("Fetching points for Telegram ID: %s", telegramID)

	var user models.User
	if err := config.DB.Where("telegram_id = ?", telegramID).First(&user).Error; err != nil {
		config.Logger.Warnf("User not found for Telegram ID: %s", telegramID)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	config.Logger.Infof("Points fetched successfully for Telegram ID: %s", telegramID)
	c.JSON(http.StatusOK, gin.H{"points": user.Points})
}

func AddPoints(c *gin.Context) {
	telegramID := c.Param("telegram_id")

	result := config.DB.Model(&models.User{}).
		Where("telegram_id = ?", telegramID).
		Update("points", gorm.Expr("points + ?", 1))

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Handle any error during the update
	if result.Error != nil {
		config.Logger.Errorf("Failed to update points for Telegram ID: %s, error: %v", telegramID, result.Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update points"})
		return
	}

	var updatedPoints int
	if err := config.DB.Model(&models.User{}).
		Where("telegram_id = ?", telegramID).
		Select("points").
		Scan(&updatedPoints).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated points"})
		return
	}

	config.Logger.Infof("Points incremented for Telegram ID: %s, New Points: %d", telegramID, updatedPoints)
	c.JSON(http.StatusOK, gin.H{"message": "Points added", "points": updatedPoints})
}
