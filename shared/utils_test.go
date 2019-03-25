package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinURL(t *testing.T) {
	tests := map[string]struct {
		args     []string
		expected string
	}{
		"no trailing slashes": {
			args:     []string{"http://example.com", "/api", "/path"},
			expected: "http://example.com/api/path",
		},
		"with trailing slashes": {
			args:     []string{"ws://example.com/", "/api/"},
			expected: "ws://example.com/api",
		},
		"no leading slashes": {
			args:     []string{"https://example.com", "api", "path"},
			expected: "https://example.com/api/path",
		},
		"no URL scheme": {
			args:     []string{"example.com", "api", "path"},
			expected: "example.com/api/path",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := JoinURL(test.args[0], test.args[1:]...)
			assert.Equal(t, test.expected, result)
		})
	}
}
