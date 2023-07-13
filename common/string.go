package common

import (
	"crypto/rand"
	"encoding/hex"
	"strings"

	"github.com/google/uuid"
)

func NormalizeString(input string) string {
	resultLower := strings.ToLower(input)
	result := strings.ReplaceAll(resultLower, " ", "_")
	return result
}

func GenerateUUID() string {
	id := uuid.New()
	result := id.String()
	return result
}
func CreateRandomFactor() (result string) {
	b := make([]byte, 4)
	rand.Read(b)
	result = hex.EncodeToString(b)
	return result
}
