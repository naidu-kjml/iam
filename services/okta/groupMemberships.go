package okta

// GroupMembership holds the current user ids for users who are part of a given group
type GroupMembership struct {
	GroupID   string
	GroupName string
	Users     []string
}

func (c *Client) fetchGroupMembership(groupID string) ([]string, error) {
	url, err := joinURL(c.baseURL, "/groups/", groupID, "/users")
	if err != nil {
		return nil, err
	}

	var allUsers []string
	var resources []struct {
		Profile oktaUserProfile
	}

	responses, err := c.fetchPagedResource(url)
	if err != nil {
		return nil, err
	}

	for _, response := range responses {
		jsonErr := response.JSON(&resources)
		if jsonErr != nil {
			return nil, jsonErr
		}

		users := make([]string, len(resources))
		for i := range resources {
			users[i] = resources[i].Profile.EmployeeNumber
		}

		allUsers = append(allUsers, users...)
	}

	return allUsers, nil
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
