package main

import (
	"log"
	"net"
	"net/http"
	"time"

	restAPI "gitlab.skypicker.com/platform/security/iam/api/rest"
	"gitlab.skypicker.com/platform/security/iam/security/permissions"
	"gitlab.skypicker.com/platform/security/iam/security/secrets"
	"gitlab.skypicker.com/platform/security/iam/services/okta"
	"gitlab.skypicker.com/platform/security/iam/storage"

	"github.com/getsentry/raven-go"
	"github.com/spf13/viper"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func fillCache(client *okta.Client) {
	// fill cache immediately (if not dev)
	if viper.GetString("APP_ENV") != "dev" {
		log.Println("Start caching users")
		client.SyncUsers()
		log.Println("Start caching groups")
		client.SyncGroups()
	}

	// Run periodic task to fill in the cache
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()

	for tick := range ticker.C {
		log.Println("Start caching users", tick.Round(time.Second))
		client.SyncUsers()
		log.Println("Start caching groups", tick.Round(time.Second))
		client.SyncGroups()
	}
}

func syncVault(client *secrets.VaultManager) {
	log.Println("Start token sync with Vault")
	err := client.SyncAppTokens()

	if err != nil {
		log.Println("[ERROR] failed to sync tokens with Vault: ", err)
	}

	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()

	for tick := range ticker.C {
		log.Println("Start caching app tokens from Vault", tick.Round(time.Second))
		err := client.SyncAppTokens()

		if err != nil {
			log.Println("[ERROR] failed to sync tokens with Vault: ", err)
		}
	}
}

func loadEnv() secrets.SecretManager {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env.yaml")
	configErr := viper.ReadInConfig()

	if configErr != nil {
		log.Println("Config file failed to load. Defaulting to env.")
	}

	// Set defaults
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("SERVE_PATH", "/")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_LOCK_RETRY_DELAY", "1s")
	viper.SetDefault("REDIS_LOCK_EXPIRATION", "5s")

	tokenConfig, err := secrets.CreateNewConfigurationMapper()
	if err != nil {
		panic(err)
	}

	// Load data from Vault and set them if possible
	vaultClient, vaultErr := secrets.CreateNewVaultClient(
		viper.GetString("VAULT_ADDR"),
		viper.GetString("VAULT_TOKEN"),
		viper.GetString("VAULT_NAMESPACE"),
		tokenConfig,
	)
	if vaultErr == nil {
		// This sync needs to happen synchronously
		err := vaultClient.SyncAppSettings()

		// If Vault is set up, but connection fails or settings are empty then kill the app as oktaToken will not be available
		if err != nil {
			panic(err)
		}

		go syncVault(vaultClient)

		return vaultClient
	}

	log.Println("Vault integration disabled: ", vaultErr)

	localSecretManager := secrets.CreateNewLocalSecretManager(tokenConfig)

	return localSecretManager
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
	secretManager := loadEnv()
	initErrorTracking(viper.GetString("SENTRY_DSN"), viper.GetString("APP_ENV"), viper.GetString("SENTRY_RELEASE"))

	permissionManager := permissions.NewYamlPermissionManager()
	go permissionManager.LoadPermissions()

	// Datadog tracer
	datadogEnv := viper.GetString("DATADOG_ENV")
	if datadogEnv != "" {
		tracer.Start(
			tracer.WithServiceName("kiwi-iam"),
			tracer.WithGlobalTag("env", datadogEnv),
		)
	}

	var port = viper.GetString("PORT")

	oktaToken, _ := secretManager.GetSetting("OKTA_TOKEN")
	cache := storage.NewRedisCache(
		viper.GetString("REDIS_HOST"),
		viper.GetString("REDIS_PORT"),
	)
	lock := storage.NewLockManager(
		cache,
		viper.GetDuration("REDIS_LOCK_RETRY_DELAY"),
		viper.GetDuration("REDIS_LOCK_EXPIRATION"),
	)
	var oktaClient = okta.NewClient(okta.ClientOpts{
		BaseURL:     viper.GetString("OKTA_URL"),
		AuthToken:   oktaToken,
		Cache:       cache,
		LockManager: lock,
	})

	router := restAPI.CreateRouter("kiwi-iam.http.router", oktaClient, permissionManager, secretManager)

	// 0.0.0.0 is specified to allow listening in Docker
	var address = "0.0.0.0"
	if viper.GetBool("USE_LOCALHOST") {
		address = "localhost"
	}

	var serveAddr = net.JoinHostPort(address, port)

	server := &http.Server{
		Handler:      router,
		Addr:         serveAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go fillCache(oktaClient)

	log.Println("🚀 Golang server starting on " + serveAddr)
	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err.Error())
	}
}
