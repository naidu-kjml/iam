package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Healthcheck : Route for healthcheck
func Healthcheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var message = "Ok"
	w.Write([]byte(message))
}
