package services

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/spf13/viper"
)

// UserProfile formatted user data
type UserProfile struct {
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
}

type apiUserProfile struct {
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

func formatUser(user apiUserProfile) UserProfile {
	return UserProfile{
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
	Profile apiUserProfile
}

// GetUserByEmail : Fetches a Okta user by email
func GetUserByEmail(email string) UserProfile {
	var oktaURL = viper.GetString("OKTA_URL")
	var oktaToken = viper.GetString("OKTA_TOKEN")

	var url = oktaURL + "/users/" + email
	log.Println("GET", url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", oktaToken)
	if err != nil {
		log.Println("Error creating new Request", err)
	}

	response, err := HTTPClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()

	var data oktaResponse
	json.NewDecoder(response.Body).Decode(&data)

	return formatUser(data.Profile)
}
