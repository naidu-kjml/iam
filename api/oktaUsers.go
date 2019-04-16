package api

import (
	"log"
	"net/http"
	"net/mail"
	"net/url"

	"github.com/getsentry/raven-go"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"gitlab.skypicker.com/platform/security/iam/services/okta"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// getOktaUserByEmail looks up an Okta user by email
func getOktaUserByEmail(client *okta.Client) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var values, err = url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		var email = values.Get("email")
		_, err = mail.ParseAddress(email)
		if err != nil {
			http.Error(w, "Invalid email", http.StatusBadRequest)
			return
		}

		userData, err := client.GetUser(email)
		if err == okta.ErrUserNotFound {
			http.Error(w, "User "+email+" not found", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Service unavailable", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		je := json.NewEncoder(w)
		if err := je.Encode(userData); err != nil {
			log.Println("[ERROR]", err.Error())
			raven.CaptureError(err, nil)
		}
	}
}
