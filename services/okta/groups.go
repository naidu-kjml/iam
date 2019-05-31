package okta

import (
	gourl "net/url"
	"strings"
	"time"
)

const iamGroupPrefix = "iam-"

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

func (c *Client) fetchGroups(userID, since string) ([]Group, error) {
	var filter string
	if since != "" {
		filter = "?filter=" + gourl.QueryEscape("lastMembershipUpdated gt \""+since+"\"")
	}

	var url string
	var err error
	if userID != "" {
		url, err = joinURL(c.baseURL, "/users/", userID, "/groups/")
	} else {
		url, err = joinURL(c.baseURL, "/groups/")
	}
	if err != nil {
		return nil, err
	}

	var allGroups []Group
	var resources []struct {
		ID                    string
		Profile               oktaGroupProfile
		LastMembershipUpdated time.Time
		LastFetched           time.Time
	}

	responses, err := c.fetchPagedResource(url + filter)
	if err != nil {
		return nil, err
	}

	for _, response := range responses {
		jsonErr := response.JSON(&resources)
		if jsonErr != nil {
			return nil, jsonErr
		}

		var groups []Group
		for i := range resources {
			group := &resources[i]
			if strings.HasPrefix(group.Profile.Name, iamGroupPrefix) {
				groups = append(groups, Group{
					ID:                    group.ID,
					Name:                  group.Profile.Name,
					Description:           group.Profile.Description,
					LastMembershipUpdated: group.LastMembershipUpdated,
				})
			}
		}
		allGroups = append(allGroups, groups...)
	}

	return allGroups, nil
}
