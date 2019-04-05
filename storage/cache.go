package storage

import (
	"net"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"gitlab.skypicker.com/platform/security/iam/shared"
	redisTrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/go-redis/redis"
)

// Cache contains cache client
type Cache struct {
	client *redisTrace.Client
	lock   *LockOpts
}

// NewCache initializes and returns a Cache
func NewCache(host, port string, lock *LockOpts) *Cache {
	opts := &redis.Options{Addr: net.JoinHostPort(host, port)}

	return &Cache{
		client: redisTrace.NewClient(opts, redisTrace.WithServiceName("kiwi-iam.redis")),
		lock:   lock,
	}
}

// Get retrieves an item from cache.
// `key` is case insensitive.
// `value` is a pointer to the variable that will receive the data.
// `error` is redis.Nil when no value is found.
func (c *Cache) Get(key string, value interface{}) error {
	lowerKey := strings.ToLower(key)
	data, err := c.client.Get(lowerKey).Bytes()
	if err != nil {
		return err
	}

	err = shared.JSON.Unmarshal(data, &value)
	return err
}

// Set writes data to cache with the specified lifespan
// `key` is case insensitive.
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) error {
	strVal, err := shared.JSON.Marshal(value)
	if err != nil {
		return err
	}

	lowerKey := strings.ToLower(key)
	_, err = c.client.Set(lowerKey, strVal, ttl).Result()
	return err
}

// Del deletes an item from cache
func (c *Cache) Del(key string) error {
	lowerKey := strings.ToLower(key)
	_, err := c.client.Del(lowerKey).Result()
	return err
}

// MSet writes items to cache in bulk
func (c *Cache) MSet(pairs map[string]interface{}) error {
	args := make([]interface{}, len(pairs)*2)
	i := 0

	for key, value := range pairs {
		args[i] = strings.ToLower(key)
		strValue, err := shared.JSON.Marshal(value)
		if err != nil {
			return err
		}
		args[i+1] = strValue
		i += 2
	}

	_, err := c.client.MSet(args...).Result()
	return err
}
