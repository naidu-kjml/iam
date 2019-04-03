package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Healthcheck : Route for healthcheck
func Healthcheck(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var message = "Ok"
	_, err := w.Write([]byte(message))

	if err != nil {
		log.Println("[ERROR] Failed to return healthcheck: ", err)
	}
}
