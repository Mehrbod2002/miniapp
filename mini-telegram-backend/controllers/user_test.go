package controllers

import (
	"bytes"
	"encoding/json"
	"log"
	"mini-telegram/config"
	"mini-telegram/models"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	var err error
	config.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to initialize in-memory database: %v", err)
	}

	// Run migrations for testing
	err = config.DB.AutoMigrate(&models.User{}, &models.LogEntry{})
	if err != nil {
		log.Fatalf("Failed to migrate test database: %v", err)
	}

	config.InitializeLogger()

	os.Exit(m.Run())
}

// Helper to reset database between tests
func resetDatabase() {
	config.DB.Exec("DELETE FROM users")
	config.DB.Exec("DELETE FROM log_entries")
}

func TestRegister(t *testing.T) {
	resetDatabase()

	router := gin.Default()
	router.POST("/register", Register)

	user := models.User{
		TelegramID: 123456,
		Username:   "testuser",
		Password:   "password",
	}

	body, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	// Verify database entry
	var dbUser models.User
	if err := config.DB.First(&dbUser, "telegram_id = ?", user.TelegramID).Error; err != nil {
		t.Errorf("Expected user to be created, got error: %v", err)
	}
}

func TestGetPoints(t *testing.T) {
	resetDatabase()

	testUser := models.User{
		TelegramID: 123456,
		Username:   "testuser",
		Points:     10,
	}
	config.DB.Create(&testUser)

	router := gin.Default()
	router.GET("/points/:telegram_id", GetPoints)

	req, _ := http.NewRequest("GET", "/points/123456", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	var response map[string]int
	if err := json.Unmarshal(resp.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if response["points"] != 10 {
		t.Errorf("Expected points to be 10, got %d", response["points"])
	}
}

func TestAddPoints(t *testing.T) {
	resetDatabase()

	testUser := models.User{
		TelegramID: 123456,
		Username:   "testuser",
		Points:     10,
	}
	config.DB.Create(&testUser)

	router := gin.Default()
	router.POST("/points/:telegram_id", AddPoints)

	req, _ := http.NewRequest("POST", "/points/123456", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.Code)
	}

	var updatedUser models.User
	config.DB.First(&updatedUser, "telegram_id = ?", 123456)
	if updatedUser.Points != 11 {
		t.Errorf("Expected points to be 11, got %d", updatedUser.Points)
	}

	// Parse and verify response body
	var response map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to parse response body: %v", err)
	}

	if response["message"] != "Points added" {
		t.Errorf("Expected message 'Points added', got %s", response["message"])
	}

	if int(response["points"].(float64)) != 11 {
		t.Errorf("Expected points to be 11, got %d", int(response["points"].(float64)))
	}
}
