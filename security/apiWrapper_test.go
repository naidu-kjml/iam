package security

import (
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCheckAuth(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/", nil)
	assert.Panics(t, func() { checkAuth(req) }, "Should panic on missing email")

	req, _ = http.NewRequest("GET", "http://example.com/?email=email@example.com", nil)
	assert.Panics(t, func() { checkAuth(req) }, "Should panic on missing User-Agent")
	req.Header.Set("User-Agent", "serviceName")

	assert.Panics(t, func() { checkAuth(req) }, "Should panic on missing Authorization header")
	req.Header.Set("Authorization", "invalid token")

	assert.Panics(t, func() { checkAuth(req) }, "Should panic on invalid token")
	req.Header.Set("Authorization", "valid token")
	viper.Set("TOKEN_serviceName_OKTA", "valid token")

	assert.NotPanics(t, func() { checkAuth(req) }, "Should not panic on valid request token")
}
