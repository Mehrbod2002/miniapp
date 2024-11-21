package controllers

import (
	"mini-telegram/config"
	"mini-telegram/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.Points = 0
	if err := config.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func GetPoints(c *gin.Context) {
	username := c.Param("username")
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"points": user.Points})
}

func AddPoints(c *gin.Context) {
	username := c.Param("username")
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Points += 10
	config.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "Points added", "points": user.Points})
}
