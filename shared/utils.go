package shared

import (
	"net/url"
	"path"

	jsoniter "github.com/json-iterator/go"
)

// JSON : faster implementation of standard JSON library
var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

// JoinURL : parses and joins a base URL to a path safely
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
