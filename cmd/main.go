package main

import (
	"log"
	"net/http"
	"time"

	"gitlab.skypicker.com/cs-devs/governant/services/okta"
	"gitlab.skypicker.com/cs-devs/governant/shared"

	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"gitlab.skypicker.com/cs-devs/governant/api"
)

func updateUserData() {

	users, err := okta.FetchAllUsers()
	if err != nil {
		log.Println("Error fetching users", err)
	}

	err = okta.CacheMSet(users)
	if err != nil {
		log.Println("Error caching users", err)
	}

	log.Println("Cached ", len(users), " users")
}

func fillCache() {

	// fill cache immediatelly (if not dev)
	if viper.GetString("APP_ENV") != "dev" {
		log.Println("Start caching users")
		updateUserData()
	}

	// Run periodic task to fill in the cache
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()

	for tick := range ticker.C {
		log.Println("Start caching users", tick.Round(time.Second))
		updateUserData()
	}
}

func panicHandler(w http.ResponseWriter, r *http.Request, err interface{}) {
	apiError, ok := err.(shared.APIError)
	if ok {
		http.Error(w, apiError.Message, apiError.Code)
		return
	}

	log.Println(apiError)
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
	router.PanicHandler = panicHandler

	go fillCache()

	err := http.ListenAndServe("localhost:"+port, router)
	if err != nil {
		panic(err)
	}
}
