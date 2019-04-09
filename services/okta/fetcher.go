package okta

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"gitlab.skypicker.com/platform/security/iam/shared"
)

type apiUser struct {
	EmployeeNumber   string
	FirstName        string
	LastName         string
	Department       string
	Email            string
	KbJobPosition    string   `json:"kb_jobPosition"`
	KbPlaceOfWork    string   `json:"kb_place_of_work"`
	KbIsVendor       bool     `json:"kb_is_vendor"`
	KbTeamMembership []string `json:"kb_team_membership"`
	Manager          string
}

func formatUser(user *apiUser) User {
	return User{
		EmployeeNumber: user.EmployeeNumber,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Department:     user.Department,
		Email:          user.Email,
		Position:       user.KbJobPosition,
		Location:       user.KbPlaceOfWork,
		IsVendor:       user.KbIsVendor,
		TeamMembership: user.KbTeamMembership,
		Manager:        user.Manager,
	}
}

type oktaResponse struct {
	Profile apiUser
}

// ErrUserNotFound is returned when a user is not present in Okta
var ErrUserNotFound = errors.New("user not found")

// fetchUser retrieves a user from Okta by email
func (c *Client) fetchUser(email string) (User, error) {
	userURL, err := shared.JoinURL(c.baseURL, "/users/", email)
	if err != nil {
		return User{}, err
	}

	var response oktaResponse
	var request = shared.Request{
		Method: "GET",
		URL:    userURL,
		Body:   nil,
		Token:  c.authToken,
	}

	httpResponse, err := shared.Fetch(request)
	if err != nil {
		return User{}, err
	}
	if httpResponse.StatusCode == http.StatusNotFound {
		return User{}, ErrUserNotFound
	}
	if httpResponse.StatusCode != http.StatusOK {
		var errorMessage = "GET " + userURL + " returned error: " + httpResponse.Status
		log.Println(errorMessage)
		return User{}, errors.New(errorMessage)
	}

	jsonErr := httpResponse.JSON(&response)

	if jsonErr != nil {
		return User{}, jsonErr
	}

	var user = formatUser(&response.Profile)
	return user, nil
}

// fetchUsers is used in iterations by FetchAllUsers
func (c *Client) fetchUsers(url string) ([]User, http.Header, error) {

	var response []oktaResponse
	var request = shared.Request{
		Method: "GET",
		URL:    url,
		Body:   nil,
		Token:  c.authToken,
	}

	httpResponse, err := shared.Fetch(request)
	if err != nil {
		return nil, nil, err
	}

	jsonErr := httpResponse.JSON(&response)

	if jsonErr != nil {
		return nil, nil, jsonErr
	}

	// Create empty slice with the same length as the response we got from Okta.
	var users = make([]User, len(response))
	for i := range response {
		user := &response[i]
		users[i] = formatUser(&user.Profile)
	}

	return users, httpResponse.Header, nil
}

// fetchAllUsers retrieves all Okta users
func (c *Client) fetchAllUsers() ([]User, error) {
	var allUsers []User

	url, err := shared.JoinURL(c.baseURL, "/users/")
	if err != nil {
		return nil, err
	}
	hasNext := true

	for hasNext {
		hasNext = false
		var users []User
		var header http.Header

		users, header, err = c.fetchUsers(url)

		if err != nil {
			return nil, err
		}

		linkHeader := header["Link"]
		for _, link := range linkHeader {
			if strings.Contains(link, "rel=\"next\"") {
				regex := regexp.MustCompile(`(?:<)(.*)(?:>)`)
				url = regex.FindStringSubmatch(link)[1]
				if url != "" {
					hasNext = true
				}
			}
		}
		allUsers = append(allUsers, users...)
	}

	return allUsers, err
}
