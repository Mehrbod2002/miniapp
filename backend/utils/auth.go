package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"mini-telegram/config"
	"net/url"
	"strings"
)

// VerifyInitData validates the initData signature provided by Telegram
func VerifyInitData(initData string) (map[string]string, bool) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, false
	}

	// Extract the hash from initData
	hash := values.Get("hash")
	if hash == "" {
		return nil, false
	}

	// Create data_check_string by sorting keys (Telegram's specification)
	var dataCheckString strings.Builder
	for key, value := range values {
		if key != "hash" {
			dataCheckString.WriteString(key + "=" + value[0] + "\n")
		}
	}
	dataCheckStringStr := strings.TrimSuffix(dataCheckString.String(), "\n")
	secret := config.GetEnv("SECRET", "")

	// Generate secret key
	secretKey := sha256.Sum256([]byte("WebAppData" + secret))
	h := hmac.New(sha256.New, secretKey[:])
	h.Write([]byte(dataCheckStringStr))
	expectedHash := hex.EncodeToString(h.Sum(nil))

	// Compare the calculated hash with the received hash
	if !hmac.Equal([]byte(expectedHash), []byte(hash)) {
		return nil, false
	}

	// Convert initData into a map for use
	data := make(map[string]string)
	for key, value := range values {
		data[key] = value[0]
	}

	return data, true
}

func OnCall(initData string) bool {
	_, valid := VerifyInitData(initData)
	return valid
}
