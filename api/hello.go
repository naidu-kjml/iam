package api

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// SayHello : Route for HelloWorld
func SayHello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
}
