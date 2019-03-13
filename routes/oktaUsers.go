package routes

import (
	"net/http"

	jsoniter "github.com/json-iterator/go"

	"github.com/julienschmidt/httprouter"
	services "gitlab.skypicker.com/cs-devs/overseer-okta/services"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// GetOktaUserByEmail : Look up Okta user by email
func GetOktaUserByEmail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var query = r.URL.Query
	var email = query()["email"][0]

	userData, err := services.GetUserByEmail(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonData)
}

// GetOktaUsers : Get all Okta users
func GetOktaUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var userData = services.GetUsers("")
	jsonData, err := json.Marshal(userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonData)
}
