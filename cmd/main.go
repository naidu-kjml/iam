package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"

	"github.com/spf13/viper"
	routes "gitlab.skypicker.com/cs-devs/overseer-okta/routes"
)

func main() {
	viper.AutomaticEnv()
	log.Println(viper.GetString("OKTA_URL"))

	var port = "8080"

	// Override for default PORT setting
	if viper.IsSet("PORT") {
		port = viper.GetString("PORT")
	}

	log.Println("Server started on " + port)
	router := httprouter.New()
	router.GET("/", routes.GetOktaUserByEmail)

	if err := http.ListenAndServe(":"+port, router); err != nil {
		panic(err)
	}
}
