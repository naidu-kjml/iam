package cache

import (
	"time"

	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"gitlab.skypicker.com/cs-devs/overseer-okta/types"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
var redisClient = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
})

// GetOkta : get a cached Okta profile. `error` is redis.Nil when no value is found.
func GetOkta(email string) (types.OktaProfile, error) {
	var profile types.OktaProfile

	data, err := redisClient.Get(email).Bytes()
	if err != nil {
		return profile, err
	}

	err = json.Unmarshal(data, &profile)
	return profile, err
}

// SetOkta : store an Okta profile to cache.
func SetOkta(key string, profile types.OktaProfile, ttl time.Duration) error {
	strProfile, err := json.Marshal(profile)
	if err != nil {
		return err
	}

	_, err = redisClient.Set(key, strProfile, ttl).Result()
	return err
}
