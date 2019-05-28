package okta

import (
	"time"

	"gitlab.skypicker.com/platform/security/iam/storage"
	"golang.org/x/sync/singleflight"
)

// ClientOpts contains options to create an Okta client
type ClientOpts struct {
	Cache       *storage.RedisCache
	LockManager *storage.LockManager
	BaseURL     string
	AuthToken   string
}

// Client represent an Okta client
type Client struct {
	group         singleflight.Group
	cache         *storage.RedisCache
	lock          *storage.LockManager
	lastGroupSync time.Time
	baseURL       string
	authToken     string
}

// NewClient creates an Okta client based on the given options
func NewClient(opts ClientOpts) *Client {
	return &Client{
		cache:         opts.Cache,
		lock:          opts.LockManager,
		baseURL:       opts.BaseURL,
		authToken:     opts.AuthToken,
		lastGroupSync: time.Unix(0, 0).UTC(),
	}
}
