package storage

import (
	"time"

	"github.com/getsentry/raven-go"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// LockOpts contains options to configure the behavior of cache locks
type LockOpts struct {
	RetryDelay time.Duration
	Expiration time.Duration
}

// ErrLockExists indicates that a lock was not created because one was already
// present. This error is returned after the old lock is deleted or expired.
var ErrLockExists = errors.New("lock was not created because one was already present")

// Lock creates a lock to prevent having multiple instances of this service
// doing an expensive action at the same time. If a lock already exists,
// the function will not create one, it will wait until the existing one is
// deleted or expired before returning ErrLockExists.
func (c *Cache) Lock(name string) error {
	var lock time.Time
	var exists bool
	key := "lock:" + name

	// Check if a lock already exists.
	err := c.Get(key, &lock)
	for err == nil {
		// If it does, wait for it to expire or be deleted.
		exists = true
		time.Sleep(c.lock.RetryDelay)
		err = c.Get(key, &lock)
	}
	if exists {
		return ErrLockExists
	}
	if err != redis.Nil {
		err = errors.Wrap(err, "error checking if a lock exists")
		raven.CaptureError(err, nil)
	}

	lock = time.Now()
	err = c.Set(key, lock, c.lock.Expiration)
	if err != nil {
		err = errors.Wrap(err, "error creating lock")
		raven.CaptureError(err, nil)
	}
	return nil
}

// Unlock removes a lock for the provided name.
func (c *Cache) Unlock(name string) {
	key := "lock:" + name
	err := c.Del(key)
	if err != nil {
		raven.CaptureError(err, nil)
	}
}
