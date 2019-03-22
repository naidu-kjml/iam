package okta

import (
	"time"

	"github.com/go-redis/redis"
)

// GetUser returns an Okta user by email. It first tries to get it from cache,
// and if not present there, it will fetch it from Okta API.
func GetUser(email string) (User, error) {
	user, err := CacheGet(email)
	if err == nil {
		// Cache hit
		return user, nil
	}

	if err != redis.Nil {
		// Not a cache hit, not a cache miss, something went wrong
		return user, err
	}

	// Cache miss
	user, err = FetchUser(email)
	if err != nil {
		return user, err
	}

	err = CacheSet(user.Email, user, time.Minute*10)
	return user, err
}
