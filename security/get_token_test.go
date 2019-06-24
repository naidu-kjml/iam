package security

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetToken(t *testing.T) {
	tests := map[string]string{
		"Bearer token":        "token",
		"token":               "token",
		"Bearer Bearer token": "Bearer token",
	}

	for test, expected := range tests {
		actual := GetToken(test)
		assert.Equal(t, expected, actual)
	}
}
