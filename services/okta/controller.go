package okta

import (
	"log"
	"time"

	"github.com/getsentry/raven-go"
	"github.com/go-redis/redis"
)

// GetUser returns an Okta user by email. It first tries to get it from cache,
// and if not present there, it will fetch it from Okta API.
func (c *Client) GetUser(email string) (User, error) {
	var user User
	err := c.cache.Get(email, &user)
	if err == nil {
		// Cache hit
		return user, nil
	}

	if err != redis.Nil {
		// Not a cache hit, not a cache miss, something went wrong
		return User{}, err
	}

	// Cache miss
	// Deduplicate network calls and cache writes if this controller is called
	// multiple times concurrently.
	val, err, _ := c.group.Do(email, func() (interface{}, error) {
		user, err := c.fetchUser(email)
		if err != nil {
			return User{}, err
		}

		err = c.cache.Set(user.Email, user, time.Minute*10)
		if err != nil {
			raven.CaptureError(err, nil)
		}
		return user, nil
	})

	if err != nil {
		return User{}, err
	}
	return val.(User), nil
}

// SyncUsers gets all users from Okta and saves them into cache.
func (c *Client) SyncUsers() {
	users, err := c.fetchAllUsers()
	if err != nil {
		log.Println("Error fetching users", err)
		raven.CaptureError(err, nil)
		return
	}

	pairs := make(map[string]interface{}, len(users))
	for _, user := range users {
		pairs[user.Email] = user
	}

	err = c.cache.MSet(pairs)
	if err != nil {
		log.Println("Error caching users", err)
		raven.CaptureError(err, nil)
		return
	}

	log.Println("Cached ", len(users), " users")
}
