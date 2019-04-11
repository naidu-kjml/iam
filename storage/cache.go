package storage

import (
	"net"
	"strings"
	"time"

	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	redisTrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// RedisCache contains redis client
type RedisCache struct {
	client *redisTrace.Client
}

// NewRedisCache initializes and returns a RedisCache
func NewRedisCache(host, port string) *RedisCache {
	opts := &redis.Options{Addr: net.JoinHostPort(host, port)}

	return &RedisCache{
		redisTrace.NewClient(opts, redisTrace.WithServiceName("kiwi-iam.redis")),
	}
}

// Get retrieves an item from cache.
// `key` is case insensitive.
// `value` is a pointer to the variable that will receive the data.
// `error` is redis.Nil when no value is found.
func (c *RedisCache) Get(key string, value interface{}) error {
	lowerKey := strings.ToLower(key)
	data, err := c.client.Get(lowerKey).Bytes()
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, &value)
	return err
}

// Set writes data to cache with the specified lifespan
// `key` is case insensitive.
func (c *RedisCache) Set(key string, value interface{}, ttl time.Duration) error {
	strVal, err := json.Marshal(value)
	if err != nil {
		return err
	}

	lowerKey := strings.ToLower(key)
	_, err = c.client.Set(lowerKey, strVal, ttl).Result()
	return err
}

// Del deletes an item from cache
func (c *RedisCache) Del(key string) error {
	lowerKey := strings.ToLower(key)
	_, err := c.client.Del(lowerKey).Result()
	return err
}

// MSet writes items to cache in bulk
func (c *RedisCache) MSet(pairs map[string]interface{}) error {
	args := make([]interface{}, len(pairs)*2)
	i := 0

	for key, value := range pairs {
		args[i] = strings.ToLower(key)
		strValue, err := json.Marshal(value)
		if err != nil {
			return err
		}
		args[i+1] = strValue
		i += 2
	}

	_, err := c.client.MSet(args...).Result()
	return err
}
