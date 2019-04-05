package api

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// SayHello is a route for HelloWorld
func SayHello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	message := "42"
	_, err := w.Write([]byte(message))
	if err != nil {
		log.Println("Hello World failed: ", err)
	}
}
