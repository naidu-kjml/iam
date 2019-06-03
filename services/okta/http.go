package okta

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/getsentry/raven-go"
	jsoniter "github.com/json-iterator/go"
	"gitlab.skypicker.com/go/packages/useragent"
	"gitlab.skypicker.com/platform/security/iam/monitoring"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

// Request contains options for an HTTP request created using shared.Fetch
type Request struct {
	Method string
	URL    string
	Token  string
	Body   io.Reader
}

// Response for an HTTP request, exposes a JSON method to get the
// retrieved data
type Response struct {
	*http.Response
}

// JSON retrieves data from HTTP response and store it in the struct pointed by
// `body` (note: `body` should be a pointer to the struct you expect the HTTP
// call to return)
func (res Response) JSON(body interface{}) error {
	// Right before returning, close the stream used for reading the
	// response's body.
	defer res.Body.Close()

	return json.NewDecoder(res.Body).Decode(&body)
}

// fetch makes an HTTP request and returns response
func (c *Client) fetch(req Request) (*Response, error) {
	log.Println(req.Method, req.URL)
	httpReq, err := http.NewRequest(req.Method, req.URL, req.Body)

	ua := useragent.UserAgent{
		Name:        "kiwi-iam",
		Environment: c.iamConfig.Environment,
		Version:     c.iamConfig.Release,
	}
	uaString, uaErr := ua.Format()
	if uaErr != nil {
		log.Println("[ERR]", uaErr)
		raven.CaptureError(uaErr, nil)
	}

	httpReq.Header.Set("User-Agent", uaString)
	httpReq.Header.Set("Authorization", req.Token)
	if err != nil {
		return nil, err
	}

	c.metrics.Incr("outgoing.requests", monitoring.Tag("url", req.URL))
	httpRes, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}

	return &Response{httpRes}, nil
}

// joinURL parses and joins a base URL to a path safely
func joinURL(baseURL string, pathname ...string) (string, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}

	// prepend u.Path to pathname slice
	elems := append([]string{u.Path}, pathname...)
	u.Path = path.Join(elems...)
	return u.String(), nil
}
