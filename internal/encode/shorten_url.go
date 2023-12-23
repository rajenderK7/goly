package encode

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateShortURL(longURL string) string {
	randomBytes := make([]byte, 16)
	rand.Read(randomBytes)

	hash := sha256.Sum256(randomBytes)
	shortURL := base64.URLEncoding.EncodeToString(hash[:6])

	return shortURL
}
