package shared

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		// scopelint has some issues here. https://github.com/kyoh86/scopelint/issues/4
		test := test

		t.Run(name, func(t *testing.T) {
			result, err := JoinURL(test.args[0], test.args[1:]...)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}
