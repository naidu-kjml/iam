package main

import (
	"log"
	"net"
	"net/http"
	"time"

	restAPI "gitlab.skypicker.com/platform/security/iam/api/rest"
	"gitlab.skypicker.com/platform/security/iam/config/cfg"
	"gitlab.skypicker.com/platform/security/iam/monitoring"
	"gitlab.skypicker.com/platform/security/iam/security/secrets"
	"gitlab.skypicker.com/platform/security/iam/services/okta"
	"gitlab.skypicker.com/platform/security/iam/storage"

	"github.com/getsentry/raven-go"
)

func fillCache(client *okta.Client) {
	log.Println("Start caching users")
	client.SyncUsers()
	log.Println("Start caching groups")
	client.SyncGroups()

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

func createSecretManager(vault cfg.VaultConfig) secrets.SecretManager {
	// Load data from Vault and set them if possible
	vaultClient, vaultErr := secrets.CreateNewVaultClient(
		vault.Address,
		vault.Token,
		vault.Namespace,
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

	localSecretManager := secrets.CreateNewLocalSecretManager()

	return localSecretManager
}

func initErrorTracking(sentry cfg.SentryConfig) {
	if sentry.Token == "" {
		log.Println("SENTRY_DSN is not set. Error logging disabled.")
		return
	}
	err := raven.SetDSN(sentry.Token)
	if err != nil {
		log.Println("[ERROR] Failed to set Raven DSN: ", err)
	}

	raven.SetEnvironment(sentry.Environment)
	raven.SetRelease(sentry.Release)
}

func main() {
	cfg.InitEnv()

	var (
		iamConfig     cfg.ServiceConfig
		oktaConfig    cfg.OktaConfig
		storageConfig cfg.StorageConfig
		datadogConfig cfg.DatadogConfig
		sentryConfig  cfg.SentryConfig
		vaultConfig   cfg.VaultConfig
	)

	// If there is an error loading the envs kill the app, as nothing will work without them.
	if err := cfg.LoadConfigs(&iamConfig, &oktaConfig, &storageConfig, &datadogConfig, &sentryConfig, &vaultConfig); err != nil {
		log.Println("[ERROR]", err.Error())
		panic(err)
	}

	initErrorTracking(sentryConfig)
	secretManager := createSecretManager(vaultConfig)

	// Datadog tracer
	tracer, _ := monitoring.CreateNewTracingService(monitoring.TracerOptions{
		ServiceName: "kiwi-iam",
		Environment: iamConfig.Environment,
		Port:        "8126",
		Host:        datadogConfig.AgentHost,
	})

	// Metrics initialization
	metricClient, metricErr := monitoring.CreateNewMetricService(monitoring.MetricSettings{
		Host:        datadogConfig.AgentHost,
		Port:        "8125",
		Namespace:   "kiwi_iam.",
		Environment: iamConfig.Environment,
	})
	if metricErr != nil {
		log.Println("[ERROR]", metricErr)
	}

	cache := storage.NewRedisCache(
		storageConfig.RedisHost,
		storageConfig.RedisPort,
	)
	lock := storage.NewLockManager(
		cache,
		storageConfig.LockRetryDelay,
		storageConfig.LockExpiration,
	)
	oktaToken, _ := secretManager.GetSetting("OKTA_TOKEN")
	oktaClient := okta.NewClient(okta.ClientOpts{
		BaseURL:     oktaConfig.URL,
		AuthToken:   oktaToken,
		Cache:       cache,
		LockManager: lock,
		IAMConfig:   &iamConfig,
		Metrics:     metricClient,
	})

	router := restAPI.CreateRouter("kiwi-iam.http.router", oktaClient, secretManager, metricClient, tracer)

	// 0.0.0.0 is specified to allow listening in Docker
	var address = "0.0.0.0"
	if iamConfig.UseLocalhost {
		address = "localhost"
	}

	var serveAddr = net.JoinHostPort(address, iamConfig.Port)
	server := &http.Server{
		Handler:      router,
		Addr:         serveAddr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if iamConfig.Environment != "dev" {
		go fillCache(oktaClient)
	}

	log.Println("🚀 Golang server starting on " + serveAddr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err.Error())
	}
}
