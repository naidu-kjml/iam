package routes

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
)

// SayHello : Route for HelloWorld
func SayHello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
}
