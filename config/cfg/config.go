package cfg

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

// InitEnv reads and set default values for environment variables.
func InitEnv() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env.yaml")

	configErr := viper.ReadInConfig()
	if configErr != nil {
		log.Println("Config file failed to load. Defaulting to env.")
	}

	// viper.Unmarshal doesn't retrieve environment variables, unless they have a
	// default value, or they are specified on .env.yaml. So to make sure all envs
	// are retrieved, we set all defaults here.
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("SERVE_PATH", "/")
	// Environment used for sentry, user agent. Removes user syncing if set to dev.
	viper.SetDefault("APP_ENV", "")
	// Uses localhost intead of 0.0.0.0, useful for OSX.
	viper.SetDefault("USE_LOCALHOST", true)

	// This is only used locally, when deployed, IAM fetches the token from Vault.
	viper.SetDefault("OKTA_TOKEN", "")
	viper.SetDefault("OKTA_URL", "")

	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_LOCK_RETRY_DELAY", "1s")
	viper.SetDefault("REDIS_LOCK_EXPIRATION", "5s")

	viper.SetDefault("SENTRY_DSN", "")
	// This value should NEVER be set manually. It's generated during docker build,
	// and it's used to track the version of the app. Useful for user agent
	// generation, and for finding regressions on Sentry.
	viper.SetDefault("SENTRY_RELEASE", "")

	viper.SetDefault("DATADOG_ENV", "")
	viper.SetDefault("DATADOG_ADDR", "")
	viper.SetDefault("DD_AGENT_HOST", "")

	viper.SetDefault("VAULT_ADDR", "")
	viper.SetDefault("VAULT_TOKEN", "")
	viper.SetDefault("VAULT_NAMESPACE", "")
}

// LoadConfigs loads environment variables into provided configStructs pointers.
//
// WARNING: configStructs should be pointers!
func LoadConfigs(configStructs ...interface{}) error {
	for _, c := range configStructs {
		err := viper.Unmarshal(c)
		if err != nil {
			return err
		}
	}

	return nil
}

// ServiceConfig stores configuration values for the IAM service.
type ServiceConfig struct {
	Port         string `mapstructure:"PORT"`
	UseLocalhost bool   `mapstructure:"USE_LOCALHOST"`
	Environment  string `mapstructure:"APP_ENV"`
	Release      string `mapstructure:"SENTRY_RELEASE"`
}

// OktaConfig stores configuration values for Okta client
type OktaConfig struct {
	URL string `mapstructure:"OKTA_URL"`
}

// StorageConfig stores configuration values for storage client.
type StorageConfig struct {
	RedisHost      string        `mapstructure:"REDIS_HOST"`
	RedisPort      string        `mapstructure:"REDIS_PORT"`
	LockRetryDelay time.Duration `mapstructure:"REDIS_LOCK_RETRY_DELAY"`
	LockExpiration time.Duration `mapstructure:"REDIS_LOCK_EXPIRATION"`
}

// DatadogConfig stores configuration values for Datadog client
type DatadogConfig struct {
	Environment string `mapstructure:"DATADOG_ENV"`
	URL         string `mapstructure:"DATADOG_ADDR"`
	AgentHost   string `mapstructure:"DD_AGENT_HOST"`
}

// SentryConfig stores configuration values for Sentry client
type SentryConfig struct {
	Token       string `mapstructure:"SENTRY_DSN"`
	Environment string `mapstructure:"APP_ENV"`
	Release     string `mapstructure:"SENTRY_RELEASE"`
}

// VaultConfig stores configuration values for Vault client
type VaultConfig struct {
	Token     string `mapstructure:"VAULT_TOKEN"`
	Address   string `mapstructure:"VAULT_ADDR"`
	Namespace string `mapstructure:"VAULT_NAMESPACE"`
}
