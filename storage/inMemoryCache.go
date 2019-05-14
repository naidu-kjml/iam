package storage

import (
	"strings"
)

// InMemoryCache is an in memory cache used as a backup when Redis is unavailable
type InMemoryCache map[string][]byte

// NewInMemoryCache initializes and returns an InMemoryCache
func NewInMemoryCache() InMemoryCache {
	return make(map[string][]byte)
}

// Get retrieves an item from cache.
// `key` is case insensitive.
// `value` is a pointer to the variable that will receive the data.
// `error` is ErrNotFound when no value is found
func (c InMemoryCache) Get(key string, value interface{}) error {
	lowerKey := strings.ToLower(key)
	data, ok := c[lowerKey]
	if ok {
		return json.Unmarshal(data, &value)
	}

	return ErrNotFound
}

// Set writes data to cache. `key` is case insensitive.
func (c InMemoryCache) Set(key string, value interface{}) error {
	strVal, err := json.Marshal(value)
	if err != nil {
		return err
	}

	lowerKey := strings.ToLower(key)
	c[lowerKey] = strVal
	return nil
}

// Del deletes an item from cache
func (c InMemoryCache) Del(key string) {
	lowerKey := strings.ToLower(key)
	delete(c, lowerKey)
}

// MSet writes items to cache in bulk
func (c InMemoryCache) MSet(pairs map[string]interface{}) error {
	bytePairs := make(map[string][]byte)

	// Go through all values and convert them to byte arrays first, then write to
	// cache. This is kept deliberately in two separate steps to ensure that if
	// there is an error nothing is written in cache.
	for key, value := range pairs {
		strValue, err := json.Marshal(value)
		if err != nil {
			return err
		}
		bytePairs[key] = strValue
	}

	for key, value := range bytePairs {
		lowerKey := strings.ToLower(key)
		c[lowerKey] = value
	}
	return nil
}
