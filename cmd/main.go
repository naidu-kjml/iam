package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	routes "gitlab.skypicker.com/cs-devs/overseer-okta/routes"
)

func main() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env.yaml")
	viper.ReadInConfig()

	viper.SetDefault("PORT", "8080")
	var port = viper.GetString("PORT")

	log.Println("OKTA_URL", viper.GetString("OKTA_URL"))
	log.Println("Server started on " + port)
	router := httprouter.New()
	router.GET("/", routes.GetOktaUserByEmail)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(err)
	}
}
