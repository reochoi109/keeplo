package idgen

import (
	"strings"

	"github.com/google/uuid"
)

func GeneratePasteID() string {
	return GenerateShortUUID(8)
}

func GenerateUUID() string {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

func GenerateShortUUID(n int) string {
	id := strings.ReplaceAll(uuid.New().String(), "-", "")
	if len(id) < n {
		return id
	}
	return id[:n]
}

func GenerateTraceID() string {
	return uuid.NewString()[:8]
}
