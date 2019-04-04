package okta

import (
	"gitlab.skypicker.com/platform/security/iam/storage"
	"golang.org/x/sync/singleflight"
)

// ClientOpts contains options to create an Okta client
type ClientOpts struct {
	CacheHost string
	CachePort string
	CacheLock *storage.LockOpts
	BaseURL   string
	AuthToken string
}

// Client represent an Okta client
type Client struct {
	group     singleflight.Group
	cache     *storage.Cache
	baseURL   string
	authToken string
}

// NewClient creates an Okta client based on the given options
func NewClient(opts ClientOpts) *Client {
	return &Client{
		cache:     storage.NewCache(opts.CacheHost, opts.CachePort, opts.CacheLock),
		baseURL:   opts.BaseURL,
		authToken: opts.AuthToken,
	}
}
