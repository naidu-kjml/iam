package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	tokens "gitlab.skypicker.com/cs-devs/overseer-auth/tokencheck"
	"net/http"
	"strconv"
	"strings"
)

func sayHello(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	message := r.URL.Path
	message = strings.TrimPrefix(message, "/")
	message = "Hello " + message
	w.Write([]byte(message))
}

type tokenVerificationBody struct {
	Token string
}

func handleVerifyToken(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	decoder := json.NewDecoder(r.Body)

	var data tokenVerificationBody
	err := decoder.Decode(&data)
	if err != nil {
		panic(err)
	}

	success := tokens.CheckToken(data.Token)

	message := `{"sucess": "` + strconv.FormatBool(success) + `"}`

	w.Header().Set("Content-type:", "application/json")

	// Set forbidden status if unsuccesful TODO: set status based on BE response
	if !success {
		w.WriteHeader(http.StatusForbidden)
	}
	w.Write([]byte(message))
}

func main() {
	router := httprouter.New()
	router.GET("/", sayHello)
	router.POST("/verify-token", handleVerifyToken)

	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}
