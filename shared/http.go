package shared

import (
	"io"
	"log"
	"net/http"
	"time"
)

var httpClient = &http.Client{
	Timeout: time.Second * 10,
}

// Request : options for an HTTP request created using shared.Fetch
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

// JSON : retrieve data from HTTP response and store it in the struct pointed by
// `body` (note: `body` should be a pointer to the struct you expect the HTTP
// call to return)
func (res Response) JSON(body interface{}) error {
	// Right before returning, close the stream used for reading the
	// response's body.
	defer res.Body.Close()

	return JSON.NewDecoder(res.Body).Decode(&body)
}

// Fetch : make an HTTP request and returns response
func Fetch(req Request) (*Response, error) {
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

	return &Response{httpRes}, nil
}
