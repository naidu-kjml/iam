package secrets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	assert.True(t, StringInSlice("some", []string{"where", "some"}))
	assert.True(t, StringInSlice("caseinsensitive", []string{"CaseInsensitive"}))
	assert.False(t, StringInSlice("some", []string{"where", "somehere"}))
}
