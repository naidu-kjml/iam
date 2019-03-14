package shared

import (
	"net/http"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var defaultTimeout = time.Second * 10

// HTTPClient : Default settings for all HTTP calls
var HTTPClient = &http.Client{
	Timeout: defaultTimeout,
}

// JSON : faster implementation of standard JSON library
var JSON = jsoniter.ConfigCompatibleWithStandardLibrary
