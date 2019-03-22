package shared

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock body for HTTP response. It implements the io.ReadWriter interface
// by exposing the Read and Close methods.
type BodyMock struct {
	mock.Mock
	Value string
}

func (b *BodyMock) Read(p []byte) (int, error) {
	// Convert string to byte array and assign it to the `p` argument
	for i, el := range []byte(b.Value) {
		p[i] = el
	}

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
	res.JSON(&expectedData)

	body.AssertNumberOfCalls(t, "Close", 1)
	assert.Equal(t, expectedData, Data{Message: "this is a test"})
}
