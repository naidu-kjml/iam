package main

import (
	"bytes"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"io/ioutil"
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

	success := CheckToken(data.Token)

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

type tokenVerificationResponse struct {
	Token string
}

// CheckToken Function to check token validity with BE. TODO: Replace with own solution in the future
func CheckToken(token string) bool {
	viper.AutomaticEnv()
	content := `{"token": "` + token + `"}`
	body := []byte(content)

	url := viper.GetString("BALKAN_BE_URL") + "/accounts/api-token-verify/"

	rs, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		panic(err) // More idiomatic way would be to print the error and die unless it's a serious error
	}

	defer rs.Body.Close()

	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		panic(err)
	}

	bodyString := string(bodyBytes)

	var response tokenVerificationResponse
	json.Unmarshal([]byte(bodyString), &response)

	return response.Token == token
}
