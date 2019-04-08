package security

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedSecretManager struct {
	mock.Mock
}

func (s *mockedSecretManager) GetAppToken(app string) (string, error) {
	if app == "SERVICENAME" {
		return "valid token", nil
	}
	return "", errors.New("wrong token bro")
}

func (s *mockedSecretManager) GetSetting(app string) (string, error) {
	if app == "SERVICENAME" {
		return "valid token", nil
	}
	return "", errors.New("wrong token bro")
}

func createFakeManager() SecretManager {
	return &mockedSecretManager{}
}

func TestGetServiceName(t *testing.T) {
	tests := []string{
		"balkan",
		"BALKAN/4704b82 (Kiwi.com sandbox)",
		"balkan/1.42.1 (Kiwi.com sandbox)",
		"balkan/1.42.1",
	}
	for _, test := range tests {
		res, err := getServiceName(test)
		assert.Equal(t, res, "BALKAN")
		assert.Equal(t, err, nil)
	}

	res, err := getServiceName("balkan-graphql/1.42.1")
	assert.Equal(t, res, "BALKAN-GRAPHQL")
	assert.Equal(t, err, nil)

	res, err = getServiceName("")
	assert.Equal(t, res, "")
	assert.Error(t, err)
}

func TestCheckAuth(t *testing.T) {
	secrets := createFakeManager()

	req, _ := http.NewRequest("GET", "http://example.com/", nil)
	assert.Error(t, checkAuth(req, secrets), "Should error on missing email")

	req, _ = http.NewRequest("GET", "http://example.com/?email=email@example.com", nil)
	assert.Error(t, checkAuth(req, secrets), "Should error on missing User-Agent")
	req.Header.Set("User-Agent", "serviceName/version (Kiwi.com environment)")

	assert.Error(t, checkAuth(req, secrets), "Should error on missing Authorization header")
	req.Header.Set("Authorization", "invalid token")

	assert.Error(t, checkAuth(req, secrets), "Should error on invalid token")
	req.Header.Set("Authorization", "valid token")

	assert.NoError(t, checkAuth(req, secrets), "Should not error on valid request token")
}
