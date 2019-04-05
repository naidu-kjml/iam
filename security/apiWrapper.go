package security

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"gitlab.skypicker.com/platform/security/iam/shared"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// AuthWrapper wraps a router to validate the authentication token
func AuthWrapper(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		err := checkAuth(r)
		if err != nil {
			if apiErr, ok := err.(shared.APIError); ok {
				http.Error(w, apiErr.Message, apiErr.Code)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}

			log.Println("[ERROR]", err.Error())
			return
		}

		// Delegate request to the given handle
		h(w, r, ps)
	}
}

// checkAuth checks if user has proper token + user agent + query fields
func checkAuth(r *http.Request) error {
	var query = r.URL.Query()
	var requestToken = r.Header.Get("Authorization")
	var service = r.Header.Get("User-Agent")

	if _, exists := query["email"]; !exists {
		return shared.APIError{Message: "Query field 'email' mandatory", Code: 401}
	}

	if service == "" {
		return shared.APIError{Message: "User-Agent header mandatory", Code: 401}
	}

	if requestToken == "" {
		return shared.APIError{Message: "Authorization header with token is mandatory", Code: 401}
	}

	var token = viper.Get("TOKEN_" + service + "_OKTA")

	if span, ok := tracer.SpanFromContext(r.Context()); ok {
		span.SetTag("  service_name", service)
	}

	if token == nil || token != requestToken {
		return shared.APIError{Message: "Incorrect token", Code: 401}
	}

	return nil
}
