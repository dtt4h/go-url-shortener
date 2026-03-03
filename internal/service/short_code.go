package service

import (
	"crypto/rand"
	"math/big"
	"strings"
)

const shortCodeLength = 8

const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateShortCode(originalURL string) string {
	var result strings.Builder
	result.Grow(shortCodeLength)

	for i := 0; i < shortCodeLength; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(base62Chars))))
		result.WriteByte(base62Chars[n.Int64()])
	}

	return result.String()
}
