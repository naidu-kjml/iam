package tokencheck

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
)

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
