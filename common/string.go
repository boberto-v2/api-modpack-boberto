package common

import (
	"strings"

	"github.com/google/uuid"
)

// "github.com/gorilla/mux"
// The name mux stands for "HTTP request multiplexer".
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
