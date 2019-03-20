package shared

import (
	"io"
	"log"
	"net/http"
	"time"
)

// Request : options for an HTTP request created using shared.Fetch
type Request struct {
	Method string
	URL    string
	Token  string
	Body   io.Reader
}

// APIError : returned on non-authorized request
type APIError struct {
	Message string
	Code    int
}

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

// Fetch : make an HTTP request and returns response
func Fetch(req Request) (*http.Response, error) {
	log.Println(req.Method, req.URL)
	httpReq, err := http.NewRequest(req.Method, req.URL, req.Body)
	httpReq.Header.Set("Authorization", req.Token)
	if err != nil {
		return nil, err
	}

	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return httpRes, nil
}

// GetRequestBody returns response on the struct pointed by
// res (note: res should be a pointer to the struct you expect the HTTP call to
// return)
func GetRequestBody(httpRes *http.Response, res interface{}) error {

	// Once the surrounding function returns (shared.Fetch) close the stream used
	// for reading the response's body.
	defer httpRes.Body.Close()

	return JSON.NewDecoder(httpRes.Body).Decode(&res)
}
