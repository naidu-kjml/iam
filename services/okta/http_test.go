package okta

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock body for HTTP response. It implements the io.ReadWriter interface
// by exposing the Read and Close methods.
type BodyMock struct {
	mock.Mock
	Value string
}

func (b *BodyMock) Read(p []byte) (int, error) {
	copy(p, b.Value)

	return len(p), nil
}

func (b *BodyMock) Close() error {
	b.Called()
	return nil
}

func TestJSON(t *testing.T) {
	var body = BodyMock{Value: `{ "message": "this is a test" }`}
	body.On("Close").Return()

	type Data struct{ Message string }
	var expectedData Data

	var res = Response{
		&http.Response{Body: &body},
	}
	err := res.JSON(&expectedData)

	if err != nil {
		panic(err)
	}

	body.AssertNumberOfCalls(t, "Close", 1)
	assert.Equal(t, expectedData, Data{Message: "this is a test"})
}

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
			result, err := joinURL(test.args[0], test.args[1:]...)
			require.NoError(t, err)
			assert.Equal(t, test.expected, result)
		})
	}
}
