package okta

import (
	"time"

	"github.com/go-redis/redis"
	"gitlab.skypicker.com/cs-devs/overseer-okta/shared"
)

var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

// CacheGet : get a cached Okta user. `error` is redis.Nil when no value is found.
func CacheGet(email string) (User, error) {
	var user User

	data, err := redisClient.Get(email).Bytes()
	if err != nil {
		return user, err
	}

	err = shared.JSON.Unmarshal(data, &user)
	return user, err
}

// CacheSet : store an Okta user to cache.
func CacheSet(key string, user User, ttl time.Duration) error {
	strUser, err := shared.JSON.Marshal(user)
	if err != nil {
		return err
	}

	_, err = redisClient.Set(key, strUser, ttl).Result()
	return err
}

// CacheMSet : cache multiple Okta users at once.
func CacheMSet(users []User) error {
	var pairs = make([]interface{}, len(users)*2)
	for i, user := range users {
		strUser, err := shared.JSON.Marshal(user)
		if err != nil {
			return err
		}

		pairs[i*2] = user.Email
		pairs[i*2+1] = strUser
	}

	_, err := redisClient.MSet(pairs...).Result()
	return err
}
