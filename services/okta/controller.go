package okta

import (
	"log"
	"strings"
	"time"

	"gitlab.skypicker.com/platform/security/iam/security/secrets"
	"gitlab.skypicker.com/platform/security/iam/storage"

	"github.com/getsentry/raven-go"
)

// User contains formatted user data provided by Okta
type User struct {
	EmployeeNumber string   `json:"employeeNumber"`
	FirstName      string   `json:"firstName"`
	LastName       string   `json:"lastName"`
	Position       string   `json:"position"`
	Department     string   `json:"department"`
	Email          string   `json:"email"`
	Location       string   `json:"location"`
	IsVendor       bool     `json:"isVendor"`
	TeamMembership []string `json:"teamMembership"`
	Manager        string   `json:"manager"`
	Permissions    []string `json:"permissions"`
}

// GetUser returns an Okta user by email. It first tries to get it from cache,
// and if not present there, it will fetch it from Okta API.
func (c *Client) GetUser(email string) (User, error) {
	var user User
	err := c.cache.Get(email, &user)
	if err == nil {
		// Cache hit
		return user, nil
	}

	if err != storage.ErrNotFound {
		// Not a cache hit, not a cache miss, something went wrong
		return User{}, err
	}

	// Cache miss
	// Deduplicate network calls and cache writes if this controller is called
	// multiple times concurrently.
	val, err, _ := c.group.Do(email, func() (interface{}, error) {
		lockErr := c.lock.Create(email)
		if lockErr == storage.ErrLockExists {
			// If there was a lock for this user, it means another instance was
			// fetching its data recently, in that case we should be able to just get
			// the data from cache.
			return c.GetUser(email)
		}
		defer c.lock.Delete(email)

		user, fetchErr := c.fetchUser(email)
		if fetchErr != nil {
			return User{}, fetchErr
		}

		cacheErr := c.cache.Set(user.Email, user, time.Minute*10)
		if cacheErr != nil {
			raven.CaptureError(cacheErr, nil)
		}
		return user, nil
	})

	if err != nil {
		return User{}, err
	}
	return val.(User), nil
}

// AddPermissions adds Okta groups to the given user object.
func (c *Client) AddPermissions(user *User, service string) error {
	cachedGroupMemberships := make(map[string][]string)

	err := c.cache.Get(service, &cachedGroupMemberships)
	if err != nil {
		if err != storage.ErrNotFound {
			return err
		}
	}

	user.Permissions = make([]string, 0)
	for groupName, users := range cachedGroupMemberships {
		if secrets.StringInSlice(user.EmployeeNumber, users) {
			user.Permissions = append(user.Permissions, groupName)
		}
	}

	return nil
}

// GetTeams retrieves from cache a map of all teams and their member count.
func (c *Client) GetTeams() (map[string]int, error) {
	var teams map[string]int
	err := c.cache.Get("teams", &teams)
	return teams, err
}

// GetGroups retrieves Okta groups
func (c *Client) GetGroups() ([]Group, error) {
	var groups []Group
	err := c.cache.Get("groups", &groups)
	return groups, err
}

// SyncUsers gets all users from Okta and saves them into cache. Also generates
// a map of all possible teams and the number of users in them, and saves them
// into cache.
func (c *Client) SyncUsers() {
	lockErr := c.lock.Create("sync_users")
	if lockErr == storage.ErrLockExists {
		log.Println("Aborted, users were already fetched")
		return
	}
	defer c.lock.Delete("sync_users")

	users, err := c.fetchAllUsers()
	if err != nil {
		log.Println("Error fetching users", err)
		raven.CaptureError(err, nil)
		return
	}

	pairs := make(map[string]interface{}, len(users))
	for i := range users {
		user := &users[i]
		pairs[user.Email] = user
	}

	err = c.cache.MSet(pairs)
	if err != nil {
		log.Println("Error caching users", err)
		raven.CaptureError(err, nil)
		return
	}
	log.Println("Cached", len(users), "users")

	nTeams, err := cacheTeams(c.cache, users)
	if err != nil {
		log.Println("Error caching teams", err)
		raven.CaptureError(err, nil)
		return
	}
	log.Println("Cached", nTeams, "teams")
}

// SyncGroups gets all groups from Okta and saves them into cache.
func (c *Client) SyncGroups() {
	syncStart := time.Now().UTC()

	groups, err := c.fetchModifiedGroups()
	if err != nil {
		log.Println("Error fetching groups", err)
		raven.CaptureError(err, nil)
		return
	}

	// We need to keep track of users assigned to various groups
	groupMemberships, err := c.fetchGroupMemberships(groups)
	if err != nil {
		log.Println("Error fetching group memberships ", err)
		raven.CaptureError(err, nil)
		return
	}

	if len(groupMemberships) > 0 {
		if err = c.updateGroupMemberships(groupMemberships); err != nil {
			log.Println("Error updating group memeberships ", err)
			raven.CaptureError(err, nil)
			return
		}

	}

	c.lastGroupSync = syncStart
	log.Println("Cached", len(groupMemberships), "group memberships")
}

func (c *Client) updateGroupMemberships(memberships []GroupMembership) error {
	cachedGroupMemberships := make(map[string][]string)

	for _, membership := range memberships {
		// iam-serviceName:rule
		groupParts := strings.SplitAfterN(membership.GroupName, ":", 2)
		serviceName := strings.Replace(strings.TrimRight(groupParts[0], ":"), iamGroupPrefix, "", 1)

		err := c.cache.Get(serviceName, &cachedGroupMemberships)
		if err != nil {
			if err != storage.ErrNotFound {
				return err
			}
		}

		cachedGroupMemberships[groupParts[1]] = membership.Users

		if err := c.cache.Set(serviceName, cachedGroupMemberships, 0); err != nil {
			return err
		}
	}

	return nil
}
