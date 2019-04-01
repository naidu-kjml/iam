package okta

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/spf13/viper"
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

func formatUser(user apiUser) User {
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

// FetchUser : Fetches a Okta user by email
func FetchUser(email string) (User, error) {
	var oktaURL = viper.GetString("OKTA_URL")
	var oktaToken = viper.GetString("OKTA_TOKEN")

	userURL, err := shared.JoinURL(oktaURL, "/users/", email)
	if err != nil {
		return User{}, err
	}

	var response oktaResponse
	var request = shared.Request{
		Method: "GET",
		URL:    userURL,
		Body:   nil,
		Token:  oktaToken,
	}

	httpResponse, err := shared.Fetch(request)
	if err != nil {
		return User{}, err
	}
	if httpResponse.StatusCode != 200 {
		var errorMessage = "GET " + userURL + " returned error: " + httpResponse.Status
		log.Println(errorMessage)
		return User{}, errors.New(errorMessage)
	}

	httpResponse.JSON(&response)

	var user = formatUser(response.Profile)
	return user, nil
}

// fetchUsers function used in iterations by FetchAllUsers
func fetchUsers(url string, token string) ([]User, http.Header, error) {

	var response []oktaResponse
	var request = shared.Request{
		Method: "GET",
		URL:    url,
		Body:   nil,
		Token:  token,
	}

	httpResponse, err := shared.Fetch(request)
	if err != nil {
		return nil, nil, err
	}

	httpResponse.JSON(&response)

	// Create empty slice with the same length as the response we got from Okta.
	var users = make([]User, len(response))
	for i, user := range response {
		users[i] = formatUser(user.Profile)
	}

	return users, httpResponse.Header, nil
}

// FetchAllUsers : Fetch all Okta users
func FetchAllUsers() ([]User, error) {

	var allUsers []User

	url, err := shared.JoinURL(viper.GetString("OKTA_URL"), "/users/")
	if err != nil {
		return nil, err
	}
	token := viper.GetString("OKTA_TOKEN")
	hasNext := true

	for hasNext {
		hasNext = false
		var users []User
		var header http.Header

		users, header, err = fetchUsers(url, token)

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
