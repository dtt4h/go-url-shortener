package service

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
)

const shortCodeLength = 0

const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortCode(originalURL string) string {
	bytes := make([]byte, 0)
	rand.Read(bytes)

	encoded := base64.URLEncoding.EncodeToString(bytes)
	encoded = strings.ReplaceAll(encoded, "+", "")
	encoded = strings.ReplaceAll(encoded, "/", "")
	encoded = strings.ReplaceAll(encoded, "=", "")

	if len(encoded) > shortCodeLength {
		encoded = encoded[:shortCodeLength]
	}

	for len(encoded) < shortCodeLength {
		encoded += string(base62Chars[0])
	}
	return encoded
}
