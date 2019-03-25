package security

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"gitlab.skypicker.com/cs-devs/governant/shared"
)

// AuthWrapper : Simple wrapper for Routers to validate token
func AuthWrapper(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		checkAuth(r)

		// Delegate request to the given handle
		h(w, r, ps)
	}
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

	var token = viper.Get("TOKEN_" + service + "_OKTA")

	if token == nil || token != requestToken {
		panic(shared.APIError{Message: "Incorrect token", Code: 401})
	}
}
