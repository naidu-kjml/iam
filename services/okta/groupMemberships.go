package okta

// GroupMembership holds the current user ids for users who are part of a given group
type GroupMembership struct {
	GroupID   string
	GroupName string
	Users     []string
}

func (c *Client) fetchGroupMembership(groupID string) ([]string, error) {
	groupUsersURL, err := joinURL(c.baseURL, "/groups/", groupID, "/users")
	if err != nil {
		return nil, err
	}

	var response []struct {
		Profile oktaUserProfile
	}
	var request = Request{
		Method: "GET",
		URL:    groupUsersURL,
		Body:   nil,
		Token:  c.authToken,
	}

	httpResponse, err := Fetch(request)
	if err != nil {
		return nil, err
	}

	jsonErr := httpResponse.JSON(&response)
	if jsonErr != nil {
		return nil, jsonErr
	}

	userIDs := make([]string, len(response))
	for i := range response {
		userIDs[i] = response[i].Profile.EmployeeNumber
	}

	return userIDs, nil
}

func (c *Client) fetchGroupMemberships(groups []Group) ([]GroupMembership, error) {
	groupMemberships := make([]GroupMembership, len(groups))

	for i, group := range groups {
		users, fetchErr := c.fetchGroupMembership(group.ID)
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
