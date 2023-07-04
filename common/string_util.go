package common

import (
	"strings"

	"github.com/gofrs/uuid"
)

func NormalizeString(input string) string {
	outLower := strings.ToLower(input)
	out := strings.ReplaceAll(outLower, " ", "_")
	return out
}

func GenerateUUID() string {
	id, _ := uuid.NewV4()
	out := id.String()
	return out
}
