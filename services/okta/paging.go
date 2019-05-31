package okta

import (
	"regexp"
	"strings"
)

var oktaLinkPattern = regexp.MustCompile(`(?:<)(.*)(?:>)`)

func (c *Client) fetchResource(url string) (*Response, error) {
	var request = Request{
		Method: "GET",
		URL:    url,
		Body:   nil,
		Token:  c.authToken,
	}

	httpResponse, err := Fetch(request)
	if err != nil {
		return &Response{}, err
	}

	return httpResponse, nil
}

func (c *Client) fetchPagedResource(url string) ([]Response, error) {
	var responses []Response
	hasNext := true

	for hasNext {
		hasNext = false

		response, err := c.fetchResource(url)
		if err != nil {
			return nil, err
		}

		linkHeader := response.Header["Link"]
		for _, link := range linkHeader {
			if strings.Contains(link, "rel=\"next\"") {
				url = oktaLinkPattern.FindStringSubmatch(link)[1]
				if url != "" {
					hasNext = true
				}
			}
		}

		responses = append(responses, *response)
	}

	return responses, nil
}
