package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"gitlab.skypicker.com/cs-devs/overseer-okta/api"
)

func main() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env.yaml")
	viper.ReadInConfig()

	viper.SetDefault("PORT", "8080")
	var port = viper.GetString("PORT")

	log.Println("OKTA_URL", viper.GetString("OKTA_URL"))
	log.Println("Server started on http://localhost:" + port)
	router := httprouter.New()
	router.GET("/", api.SayHello)
	router.GET("/user/okta", api.GetOktaUserByEmail)

	if err := http.ListenAndServe("localhost:"+port, router); err != nil {
		panic(err)
	}
}
