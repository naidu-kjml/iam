package services

import (
	"net/http"
	"time"
)

var defaultTimeout = time.Second * 10

// HTTPClient : Default settings for all HTTP calls
var HTTPClient = &http.Client{
	Timeout: defaultTimeout,
}
