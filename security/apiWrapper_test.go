package security

import (
	"net/http"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCheckAuth(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/", nil)
	assert.Error(t, checkAuth(req), "Should error on missing email")

	req, _ = http.NewRequest("GET", "http://example.com/?email=email@example.com", nil)
	assert.Error(t, checkAuth(req), "Should error on missing User-Agent")
	req.Header.Set("User-Agent", "serviceName")

	assert.Error(t, checkAuth(req), "Should error on missing Authorization header")
	req.Header.Set("Authorization", "invalid token")

	assert.Error(t, checkAuth(req), "Should error on invalid token")
	req.Header.Set("Authorization", "valid token")
	viper.Set("TOKEN_serviceName_OKTA", "valid token")

	assert.NoError(t, checkAuth(req), "Should not error on valid request token")
}
