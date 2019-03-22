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
	viper.SetDefault("SERVE_PATH", "/")

	ravenDSN := viper.GetString("SENTRY_DSN")
	if ravenDSN != "" {
		raven.SetDSN(ravenDSN)
	} else {
		log.Println("SENTRY_DSN is not set. Error logging disabled.")
	}

	okta.InitCache()
}

func main() {
	var port = viper.GetString("PORT")

	// For deployments where we're not on root
	var servePath = viper.GetString("SERVE_PATH")

	router := httprouter.New()

	// Healthcheck routes. Exposed on both /healthcheck and /servePath/healthcheck to allow easier k8s set up
	router.GET("/healthcheck", api.Healthcheck)

	// Prevent setting two routes
	if servePath != "/" {
		router.GET(servePath+"healthcheck", api.Healthcheck)
	}

	// App Routes
	router.GET(servePath, api.SayHello)
	router.GET(servePath+"user/okta", api.GetOktaUserByEmail)
	router.PanicHandler = panicHandler

	// 0.0.0.0 is specified to allow listening in Docker
	var address = "0.0.0.0:" + port
	server := &http.Server{
		Handler:      router,
		Addr:         address,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go fillCache()

	log.Println("🚀 Golang server starting on " + address + servePath)
	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}
