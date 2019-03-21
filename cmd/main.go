package main

import (
	"log"
	"net/http"
	"time"

	"gitlab.skypicker.com/cs-devs/governant/services/okta"
	"gitlab.skypicker.com/cs-devs/governant/shared"

	"github.com/getsentry/raven-go"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"gitlab.skypicker.com/cs-devs/governant/api"
)

func updateUserData() {

	users, err := okta.FetchAllUsers()
	if err != nil {
		log.Println("Error fetching users", err)
		raven.CaptureError(err, nil)
		return
	}

	err = okta.CacheMSet(users)
	if err != nil {
		log.Println("Error caching users", err)
		raven.CaptureError(err, nil)
		return
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
	apiError, castingOk := err.(shared.APIError)
	if castingOk {
		raven.CaptureError(apiError, nil)
		http.Error(w, apiError.Message, apiError.Code)
		log.Println("[ERROR]", apiError)
		return
	}

	if errorType, castingOk := err.(error); castingOk {
		raven.CaptureError(errorType, nil)
	}
	log.Panic(err)
}

// Triggered before main()
func init() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env.yaml")
	viper.ReadInConfig()
	viper.SetDefault("PORT", "8080")

	ravenDSN := viper.GetString("SENTRY_DSN")
	if ravenDSN != "" {
		raven.SetDSN(ravenDSN)
	}
}

func main() {
	var port = viper.GetString("PORT")

	router := httprouter.New()
	router.GET("/", api.SayHello)
	router.GET("/user/okta", api.GetOktaUserByEmail)
	router.PanicHandler = panicHandler

	go fillCache()

	log.Println("🚀 Golang server starting on http://localhost:" + port)
	err := http.ListenAndServe("localhost:"+port, router)
	if err != nil {
		panic(err)
	}
}
