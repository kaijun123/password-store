package util

import (
	"crypto/sha256"
	"math/rand"
	"time"
)

// Reference: https://www.educative.io/answers/how-to-generate-a-random-string-of-fixed-length-in-go
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	result := make([]byte, length)
	for i := range result {
		randomNumber := random.Intn(10000)
		index := randomNumber % len(charset)
		result[i] = charset[index]
	}
	return string(result)
}

func Hash(input []byte) []byte {
	h := sha256.New()
	h.Write(input)
	hash := h.Sum(nil)
	return hash
}
