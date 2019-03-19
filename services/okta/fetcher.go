package okta

import (
	"github.com/spf13/viper"
	"gitlab.skypicker.com/cs-devs/governant/shared"
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

	var response oktaResponse
	var request = shared.Request{
		Method: "GET",
		URL:    oktaURL + "/users/" + email,
		Body:   nil,
		Token:  oktaToken,
	}

	err := shared.Fetch(request, &response)
	if err != nil {
		return User{}, err
	}

	var user = formatUser(response.Profile)
	return user, nil
}

// FetchUsers : Fetch all Okta users
func FetchUsers(after string) ([]User, error) {
	var oktaURL = viper.GetString("OKTA_URL")
	var oktaToken = viper.GetString("OKTA_TOKEN")

	var response []oktaResponse
	var request = shared.Request{
		Method: "GET",
		URL:    oktaURL + "/users/?after=" + after,
		Body:   nil,
		Token:  oktaToken,
	}

	err := shared.Fetch(request, &response)
	if err != nil {
		return nil, err
	}

	// Create empty slice with the same length as the response we got from Okta.
	var users = make([]User, len(response))
	for i, user := range response {
		users[i] = formatUser(user.Profile)
	}
	return users, nil

	// TODO get after parameter from the header below and make more requests to
	// get the rest of the users.
	// Link <https://kiwi.oktapreview.com/api/v1/users?after=000uiq5gshbbBhVnDO0h7&limit=200>; rel="next"
}
