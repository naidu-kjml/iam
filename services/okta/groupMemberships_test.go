package okta

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func (c *MockCache) Get(key string, value interface{}) error {
	return nil
}
func (c *MockCache) Set(key string, value interface{}, ttl time.Duration) error {
	return nil
}
func (c *MockCache) Del(key string) error {
	return nil
}
func (c *MockCache) MSet(pairs map[string]interface{}, ttl time.Duration) error {
	return nil
}

func TestGroupMemberships(t *testing.T) {
	cache := &MockCache{}
	client := NewClient(ClientOpts{
		Cache: cache,
	})

	tests := []struct {
		groupName   string
		shouldError bool
	}{
		{"iam-service.valid.permission", false},
		{"iam-long-service-name.with-long.permission", false},
		{"iam-service:invalid:permission", true},
		{"service", true},
	}

	for _, test := range tests {
		err := client.updateGroupMemberships([]GroupMembership{
			{"group-id", test.groupName, []string{"user1", "user2"}},
		})

		if test.shouldError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
