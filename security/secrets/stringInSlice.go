package secrets

import "strings"

// StringInSlice checks whether str matches (case-insensitive) any string in slice.
func StringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, str) {
			return true
		}
	}
	return false
}
