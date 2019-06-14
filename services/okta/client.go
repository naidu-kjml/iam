package okta

import (
	"time"

	cfg "gitlab.skypicker.com/platform/security/iam/config"
	"gitlab.skypicker.com/platform/security/iam/monitoring"
	"gitlab.skypicker.com/platform/security/iam/storage"
	"golang.org/x/sync/singleflight"
)

// Cacher contains methods needed from a cache
type Cacher interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}, ttl time.Duration) error
	Del(key string) error
	MSet(pairs map[string]interface{}, ttl time.Duration) error
}

// ClientOpts contains options to create an Okta client
type ClientOpts struct {
	Cache       Cacher
	LockManager *storage.LockManager
	BaseURL     string
	AuthToken   string
	IAMConfig   *cfg.ServiceConfig
	Metrics     *monitoring.Metrics
}

// Client represent an Okta client
type Client struct {
	group     singleflight.Group
	cache     Cacher
	lock      *storage.LockManager
	baseURL   string
	authToken string
	iamConfig *cfg.ServiceConfig
	metrics   *monitoring.Metrics
}

// NewClient creates an Okta client based on the given options
func NewClient(opts ClientOpts) *Client {
	return &Client{
		cache:     opts.Cache,
		lock:      opts.LockManager,
		baseURL:   opts.BaseURL,
		authToken: opts.AuthToken,
		iamConfig: opts.IAMConfig,
		metrics:   opts.Metrics,
	}
}
