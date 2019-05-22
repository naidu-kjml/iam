package okta

import (
	"net/http"
	"strings"
	"time"
)

// Group represents an Okta group
type Group struct {
	ID                    string    `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	LastMembershipUpdated time.Time `json:"lastMembershipUpdated"`
}

type oktaGroupProfile struct {
	Name        string
	Description string
}

func (c *Client) fetchGroups(url string) ([]Group, http.Header, error) {
	var response []struct {
		ID                    string
		Profile               oktaGroupProfile
		LastMembershipUpdated time.Time
		LastFetched           time.Time
	}
	var request = Request{
		Method: "GET",
		URL:    url,
		Body:   nil,
		Token:  c.authToken,
	}

	httpResponse, err := Fetch(request)
	if err != nil {
		return nil, nil, err
	}

	jsonErr := httpResponse.JSON(&response)
	if jsonErr != nil {
		return nil, nil, jsonErr
	}

	// Create empty slice with the same length as the response
	var groups = make([]Group, len(response))
	for i := range response {
		group := &response[i]
		groups[i] = Group{
			ID:                    group.ID,
			Name:                  group.Profile.Name,
			Description:           group.Profile.Description,
			LastMembershipUpdated: group.LastMembershipUpdated,
		}
	}
	return groups, httpResponse.Header, nil
}

func (c *Client) fetchAllGroups() ([]Group, error) {
	var allGroups []Group

	url, err := joinURL(c.baseURL, "/groups/")
	if err != nil {
		return nil, err
	}
	hasNext := true

	for hasNext {
		hasNext = false
		groups, header, fetchErr := c.fetchGroups(url)
		if fetchErr != nil {
			return nil, fetchErr
		}

		linkHeader := header["Link"]
		for _, link := range linkHeader {
			if strings.Contains(link, `rel="next"`) {
				url = oktaLinkPattern.FindStringSubmatch(link)[1]
				if url != "" {
					hasNext = true
				}
			}
		}
		allGroups = append(allGroups, groups...)
	}

	return allGroups, nil
}
