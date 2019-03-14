package okta

import (
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"gitlab.skypicker.com/cs-devs/overseer-okta/shared"
)

type apiProfile struct {
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

func formatProfile(user apiProfile) Profile {
	return Profile{
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
	Profile apiProfile
}

// GetUserByEmail : Fetches a Okta user by email
func GetUserByEmail(email string) (Profile, error) {
	user, err := CacheGet(email)
	if err == nil {
		// Cache hit
		return user, nil
	}
	if err != redis.Nil {
		// Error retrieving item
		return user, err
	}

	// Cache miss
	var oktaURL = viper.GetString("OKTA_URL")
	var oktaToken = viper.GetString("OKTA_TOKEN")

	var url = oktaURL + "/users/" + email
	log.Println("GET", url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", oktaToken)
	if err != nil {
		return user, err
	}

	response, err := shared.HTTPClient.Do(req)
	if err != nil {
		return user, err
	}
	defer response.Body.Close()

	var data oktaResponse
	shared.JSON.NewDecoder(response.Body).Decode(&data)

	user = formatProfile(data.Profile)
	err = CacheSet(user.Email, user, time.Minute)
	return user, err
}

// GetUsers : Fetch all Okta users
func GetUsers(after string) []Profile {
	var oktaURL = viper.GetString("OKTA_URL")
	var oktaToken = viper.GetString("OKTA_TOKEN")

	var url = oktaURL + "/users/?after=" + after
	log.Println("GET", url)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", oktaToken)
	if err != nil {
		log.Println("Error creating new Request", err)
	}

	response, err := shared.HTTPClient.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer response.Body.Close()

	var data []oktaResponse
	shared.JSON.NewDecoder(response.Body).Decode(&data)

	var users = make([]Profile, len(data))
	for i, user := range data {
		users[i] = formatProfile(user.Profile)
	}
	return users

	// TODO get after parameter from the header below and make more requests to
	// get the rest of the users.
	// Link <https://kiwi.oktapreview.com/api/v1/users?after=000uiq5gshbbBhVnDO0h7&limit=200>; rel="next"
}
