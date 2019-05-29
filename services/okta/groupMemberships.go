package okta

import (
	"net/http"
	"strings"
)

// GroupMembership holds the current user ids for users who are part of a given group
type GroupMembership struct {
	GroupID   string
	GroupName string
	Users     []string
}

func (c *Client) fetchGroupMembership(url string) ([]string, http.Header, error) {
	var response []struct {
		Profile oktaUserProfile
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

	userIDs := make([]string, len(response))
	for i := range response {
		userIDs[i] = response[i].Profile.EmployeeNumber
	}

	return userIDs, httpResponse.Header, nil
}

func (c *Client) fetchCompleteGroupMembership(groupID string) ([]string, error) {
	var allUsers []string

	url, err := joinURL(c.baseURL, "/groups/", groupID, "/users")
	if err != nil {
		return nil, err
	}
	hasNext := true

	for hasNext {
		hasNext = false
		var users []string
		var header http.Header

		users, header, err = c.fetchGroupMembership(url)
		if err != nil {
			return nil, err
		}

		linkHeader := header["Link"]
		for _, link := range linkHeader {
			if strings.Contains(link, "rel=\"next\"") {
				url = oktaLinkPattern.FindStringSubmatch(link)[1]
				if url != "" {
					hasNext = true
				}
			}
		}
		allUsers = append(allUsers, users...)
	}

	return allUsers, err
}

func (c *Client) fetchGroupMemberships(groups []Group) ([]GroupMembership, error) {
	groupMemberships := make([]GroupMembership, len(groups))

	for i, group := range groups {
		users, fetchErr := c.fetchCompleteGroupMembership(group.ID)
		if fetchErr != nil {
			return nil, fetchErr
		}
		groupMemberships[i] = GroupMembership{
			group.ID,
			group.Name,
			users,
		}
	}

	return groupMemberships, nil
}
