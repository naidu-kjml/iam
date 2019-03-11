package routes

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	services "gitlab.skypicker.com/cs-devs/overseer-okta/services"
)

// GetOktaUserByEmail : Look up Okta user by email
func GetOktaUserByEmail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var query = r.URL.Query
	var email = query()["email"][0]

	var userData = services.GetUserByEmail(email)
	jsonData, err := json.Marshal(userData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Write(jsonData)
}
