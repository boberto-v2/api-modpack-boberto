package common

import (
	"strings"

	"github.com/gofrs/uuid"
)

func NormalizeString(input string) string {
	resultLower := strings.ToLower(input)
	result := strings.ReplaceAll(resultLower, " ", "_")
	return result
}

func GenerateUUID() string {
	id, _ := uuid.NewV4()
	result := id.String()
	return result
}
