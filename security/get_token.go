package security

import (
	"strings"
)

// GetToken extracts token from Bearer scheme
func GetToken(authorization string) string {
	return strings.Replace(authorization, "Bearer ", "", 1)
}
