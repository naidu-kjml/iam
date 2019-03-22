package api

import (
	"net/http"

	"gitlab.skypicker.com/cs-devs/governant/shared"

	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"gitlab.skypicker.com/cs-devs/governant/services/okta"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// GetOktaUserByEmail : Look up Okta user by email
func GetOktaUserByEmail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var query = r.URL.Query

	checkAuth(r)

	var email = query()["email"][0]

	userData, err := okta.GetUser(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

// checkAuth : checks if user has proper token + user agent + query fields
func checkAuth(r *http.Request) {

	var query = r.URL.Query()
	var requestToken = r.Header.Get("Authorization")
	var service = r.Header.Get("User-Agent")

	if _, exists := query["email"]; !exists {
		panic(shared.APIError{Message: "Query field 'email' mandatory", Code: 401})
	}

	if service == "" {
		panic(shared.APIError{Message: "User-Agent header mandatory", Code: 401})
	}

	if requestToken == "" {
		panic(shared.APIError{Message: "Authorization header with token is mandatory", Code: 401})
	}

	var token = viper.Get("TOKEN." + service + ".OKTA")

	if token == nil || token != requestToken {
		panic(shared.APIError{Message: "Incorrect token", Code: 401})
	}
}
