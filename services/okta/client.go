package okta

import (
	cfg "gitlab.skypicker.com/platform/security/iam/config"
	"gitlab.skypicker.com/platform/security/iam/monitoring"
	"gitlab.skypicker.com/platform/security/iam/storage"
	"golang.org/x/sync/singleflight"
)

// ClientOpts contains options to create an Okta client
type ClientOpts struct {
	Cache       *storage.RedisCache
	LockManager *storage.LockManager
	BaseURL     string
	AuthToken   string
	IAMConfig   *cfg.ServiceConfig
	Metrics     *monitoring.Metrics
}

// Client represent an Okta client
type Client struct {
	group     singleflight.Group
	cache     *storage.RedisCache
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
