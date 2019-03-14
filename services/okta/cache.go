package okta

import (
	"time"

	"github.com/go-redis/redis"
	"gitlab.skypicker.com/cs-devs/overseer-okta/shared"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

// CacheGet : get a cached Okta profile. `error` is redis.Nil when no value is found.
func CacheGet(email string) (Profile, error) {
	var profile Profile

	data, err := redisClient.Get(email).Bytes()
	if err != nil {
		return profile, err
	}

	err = shared.JSON.Unmarshal(data, &profile)
	return profile, err
}

// CacheSet : store an Okta profile to cache.
func CacheSet(key string, profile Profile, ttl time.Duration) error {
	strProfile, err := shared.JSON.Marshal(profile)
	if err != nil {
		return err
	}

	_, err = redisClient.Set(key, strProfile, ttl).Result()
	return err
}
