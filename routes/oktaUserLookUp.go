package routes

import (
	"github.com/julienschmidt/httprouter"
	services "gitlab.skypicker.com/cs-devs/overseer-okta/services"
	"net/http"
	// "strings"
)

// GetOktaUserByEmail : Look up Okta user by email
func GetOktaUserByEmail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	query := r.URL.Query
	var email = query()["email"][0]

	services.GetUserByEmail(email)

	// response, err := services.HTTPClient.Get(url)
	w.Write([]byte("hi"))
}
