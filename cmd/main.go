package main

import (
	"log"
	"net/http"
	"time"

	"gitlab.skypicker.com/cs-devs/overseer-okta/services/okta"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"gitlab.skypicker.com/cs-devs/overseer-okta/api"
)

func fillCache(ticker *time.Ticker) {
	for tick := range ticker.C {
		log.Println("Start caching users", tick.Round(time.Second))
		users, err := okta.FetchUsers("")
		if err != nil {
			log.Println("Error fetching users", err)
		}

		err = okta.CacheMSet(users)
		if err != nil {
			log.Println("Error caching users", err)
		}
	}
}

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

	// Run periodic task to fill in the cache
	ticker := time.NewTicker(time.Minute)
	go fillCache(ticker)
	defer ticker.Stop()

	err := http.ListenAndServe("localhost:"+port, router)
	if err != nil {
		panic(err)
	}
}
