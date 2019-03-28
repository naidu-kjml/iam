package main

import (
	"log"
	"net/http"
	"time"

	"gitlab.skypicker.com/cs-devs/governant/api"
	"gitlab.skypicker.com/cs-devs/governant/security"
	"gitlab.skypicker.com/cs-devs/governant/services/okta"
	"gitlab.skypicker.com/cs-devs/governant/shared"
	"gitlab.skypicker.com/cs-devs/governant/storage"

	"github.com/getsentry/raven-go"
	"github.com/spf13/viper"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/julienschmidt/httprouter"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func updateUserData(cache *storage.Cache) {
	users, err := okta.FetchAllUsers()
	if err != nil {
		log.Println("Error fetching users", err)
		raven.CaptureError(err, nil)
		return
	}

	pairs := make(map[string]interface{}, len(users))
	for _, user := range users {
		pairs[user.Email] = user
	}

	err = cache.MSet(pairs)
	if err != nil {
		log.Println("Error caching users", err)
		raven.CaptureError(err, nil)
		return
	}

	log.Println("Cached ", len(users), " users")
}

func fillCache() {
	cache := storage.NewCache(
		viper.GetString("REDIS_HOST"),
		viper.GetString("REDIS_PORT"),
	)

	// fill cache immediately (if not dev)
	if viper.GetString("APP_ENV") != "dev" {
		log.Println("Start caching users")
		updateUserData(cache)
	}

	// Run periodic task to fill in the cache
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()

	for tick := range ticker.C {
		log.Println("Start caching users", tick.Round(time.Second))
		updateUserData(cache)
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
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")

	ravenDSN := viper.GetString("SENTRY_DSN")
	if ravenDSN != "" {
		raven.SetDSN(ravenDSN)
		raven.SetEnvironment(viper.GetString("APP_ENV"))
		raven.SetRelease(viper.GetString("SENTRY_RELEASE"))
	} else {
		log.Println("SENTRY_DSN is not set. Error logging disabled.")
	}

	// Datadog tracer
	datadogEnv := viper.GetString("DATADOG_ENV")
	if datadogEnv != "" {
		tracer.Start(
			tracer.WithServiceName("governant"),
			tracer.WithGlobalTag("env", datadogEnv),
		)
	}
}

func main() {
	var port = viper.GetString("PORT")
	cache := storage.NewCache(
		viper.GetString("REDIS_HOST"),
		viper.GetString("REDIS_PORT"),
	)

	// For deployments where we're not on root
	var servePath = viper.GetString("SERVE_PATH")

	router := httprouter.New(httprouter.WithServiceName("governant.http.router"))

	// Healthcheck routes. Exposed on both /healthcheck and /servePath/healthcheck to allow easier k8s set up
	router.GET("/healthcheck", api.Healthcheck)

	// Prevent setting two routes
	if servePath != "/" {
		router.GET(servePath+"healthcheck", api.Healthcheck)
	}

	// App Routes
	router.GET(servePath, api.SayHello)
	router.GET(servePath+"user/okta", security.AuthWrapper(api.GetOktaUserByEmail(cache)))

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
		log.Fatal(err.Error())
	}
}
