package api

import (
	"log"
	"net/http"
	"net/mail"
	"net/url"

	"github.com/getsentry/raven-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"gitlab.skypicker.com/platform/security/iam/security"
	"gitlab.skypicker.com/platform/security/iam/services/okta"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type userDataService interface {
	GetUser(string) (okta.User, error)
}

type permissionManager interface {
	GetUserPermissions(string, []string) ([]string, error)
}

// getOktaUserByEmail looks up an Okta user by email
func getOktaUserByEmail(client userDataService, permissionManager permissionManager) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var values, err = url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		var email = values.Get("email")
		if email == "" {
			http.Error(w, "Missing email", http.StatusBadRequest)
			return
		}
		_, err = mail.ParseAddress(email)
		if err != nil {
			http.Error(w, "Invalid email", http.StatusBadRequest)
			return
		}

		service, err := security.GetServiceName(r.Header.Get("User-Agent"))
		if err != nil {
			http.Error(w, "Invalid user agent", http.StatusBadRequest)
			return
		}

		oktaUser, err := client.GetUser(email)
		if err == okta.ErrUserNotFound {
			http.Error(w, "User "+email+" not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusInternalServerError)
			return
		}

		permissions, err := permissionManager.GetUserPermissions(service, oktaUser.TeamMembership)
		if err != nil {
			log.Println("[ERROR]", err.Error())
			raven.CaptureError(err, nil)
		}
		oktaUser.Permissions = permissions

		w.Header().Set("Content-Type", "application/json")
		je := json.NewEncoder(w)

		if err := je.Encode(oktaUser); err != nil {
			log.Println("[ERROR]", err.Error())
			raven.CaptureError(err, nil)
		}
	}
}
