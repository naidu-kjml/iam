package api

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"gitlab.skypicker.com/cs-devs/governant/services/okta"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// GetOktaUserByEmail : Look up Okta user by email
func GetOktaUserByEmail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var query = r.URL.Query

	var email = query()["email"][0]

	userData, err := okta.GetUser(email)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(userData)
	if err != nil {
		http.Error(w, "Service unavailable", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}
