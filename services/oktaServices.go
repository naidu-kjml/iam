package services

import (
	"encoding/json"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

type userProfile struct {
	EmployeeNumber string
	FirstName      string
	LastName       string
	Position       string
	Department     string
	Email          string
	Location       string
	isVendor       bool
	TeamMembership []string
	Manager        string
}

type apiUserProfile struct {
	EmployeeNumber   string
	FirstName        string
	LastName         string
	Department       string
	Email            string
	KbJobPosition    string   `json:"kb_job_position"`
	KbPlaceOfWork    string   `json:"kb_place_of_work"`
	KbIsVendor       bool     `json:"kb_is_vendor"`
	KbTeamMembership []string `json:"kb_team_membership"`
	Manager          string
}

type oktaResponse struct {
	Profile apiUserProfile
}

// GetUserByEmail : Fetches a Okta user by email
func GetUserByEmail(email string) {
	var url = viper.GetString("OKTA_URL") + "/users?email=" + email
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", viper.GetString("OKTA_TOKEN"))
	log.Println(req.Header)
	if err != nil {
		log.Println(err)
	}

	log.Println(url)
	response, err := HTTPClient.Do(req)

	if err != nil {
		log.Println(err)
	}

	log.Println(response)

	defer response.Body.Close()

	var result []oktaResponse

	json.NewDecoder(response.Body).Decode(&result)

	log.Println(result)
}
