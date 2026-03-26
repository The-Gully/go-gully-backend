package google

import (
	"crypto/rand"
	"encoding/base64"
	"log"
)

func GenerateRandomString(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		log.Printf("[AUTH] Failed to generate random string: %v", err)
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)[:length]
}
