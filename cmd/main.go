package main

import (
	"log"
	"net/http"
	"time"

	"gitlab.skypicker.com/platform/security/iam/api"
	"gitlab.skypicker.com/platform/security/iam/security"
	"gitlab.skypicker.com/platform/security/iam/services/okta"

	"github.com/getsentry/raven-go"
	"github.com/spf13/viper"
	"gopkg.in/DataDog/dd-trace-go.v1/contrib/julienschmidt/httprouter"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func fillCache(client *okta.Client) {
	// fill cache immediately (if not dev)
	if viper.GetString("APP_ENV") != "dev" {
		log.Println("Start caching users")
		client.SyncUsers()
	}

	// Run periodic task to fill in the cache
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()

	for tick := range ticker.C {
		log.Println("Start caching users", tick.Round(time.Second))
		client.SyncUsers()
	}
}

func loadEnv() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env.yaml")
	err := viper.ReadInConfig()

	if err != nil {
		log.Println("Config file failed to load. Defaulting to env.")
	}

	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("SERVE_PATH", "/")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
}

func initErrorTracking(token, environment, release string) {
	if token == "" {
		log.Println("SENTRY_DSN is not set. Error logging disabled.")
		return
	}
	err := raven.SetDSN(token)
	if err != nil {
		log.Println("[ERROR] Failed to set Raven DSN: ", err)
	}

	raven.SetEnvironment(environment)
	raven.SetRelease(release)
}

func main() {
	loadEnv()

	initErrorTracking(viper.GetString("SENTRY_DSN"), viper.GetString("APP_ENV"), viper.GetString("SENTRY_RELEASE"))

	// Datadog tracer
	datadogEnv := viper.GetString("DATADOG_ENV")
	if datadogEnv != "" {
		tracer.Start(
			tracer.WithServiceName("governant"),
			tracer.WithGlobalTag("env", datadogEnv),
		)
	}

	// For deployments where we're not on root
	var servePath = viper.GetString("SERVE_PATH")
	var port = viper.GetString("PORT")

	var oktaClient = okta.NewClient(okta.ClientOpts{
		BaseURL:   viper.GetString("OKTA_URL"),
		AuthToken: viper.GetString("OKTA_TOKEN"),
		CacheHost: viper.GetString("REDIS_HOST"),
		CachePort: viper.GetString("REDIS_PORT"),
	})

	router := httprouter.New(httprouter.WithServiceName("governant.http.router"))

	// Healthcheck routes. Exposed on both /healthcheck and /servePath/healthcheck to allow easier k8s set up
	router.GET("/healthcheck", api.Healthcheck)

	// Prevent setting two routes
	if servePath != "/" {
		router.GET(servePath+"healthcheck", api.Healthcheck)
	}

	// App Routes
	router.GET(servePath, api.SayHello)
	router.GET(servePath+"user/okta", security.AuthWrapper(api.GetOktaUserByEmail(oktaClient)))

	router.PanicHandler = api.PanicHandler

	// 0.0.0.0 is specified to allow listening in Docker
	var address = "0.0.0.0:" + port
	server := &http.Server{
		Handler:      router,
		Addr:         address,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go fillCache(oktaClient)

	log.Println("🚀 Golang server starting on " + address + servePath)
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err.Error())
	}
}
