package shared

import (
	"net/url"
	"path"
)

// JoinURL parses and joins a base URL to a path safely
func JoinURL(baseURL string, pathname ...string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// prepend u.Path to pathname slice
	elems := append([]string{u.Path}, pathname...)
	u.Path = path.Join(elems...)
	return u.String(), nil
}

// StringInSlice checks whether str matches any string in slice.
func StringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
