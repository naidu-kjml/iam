package security

// This package expects the following structure for Vault.
// /secret/governant/app_tokens
// /secret/governant/settings
// @TODO: Migrate to IAM naming

import (
	"errors"
	"fmt"
	"log"

	vault "github.com/hashicorp/vault/api"
)

type localStorage struct {
	// Tokens for apps integrating with IAM
	tokens map[string]string

	// Settings for IAM
	settings map[string]string
}

// SecretManager is the local struct for Vault connection
type SecretManager struct {
	client  *vault.Client
	storage localStorage
}

var namespace = "/secret/governant"

// CreateNewSecretManager create a new client to connect to Vault
func CreateNewSecretManager(address, token string) *SecretManager {
	client, _ := vault.NewClient(&vault.Config{
		Address: address,
	})

	client.SetToken(token)

	client.SetNamespace(namespace)

	return &SecretManager{client: client}
}

// SyncAppTokens syncs all the available tokens from Vault and saves them to local state
func (s *SecretManager) SyncAppTokens() error {
	data, err := s.fetchData("app_tokens")

	if err != nil {
		return err
	}

	mappedTokens := make(map[string]string, len(data))

	for key, value := range data {
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid conversion to string for value <%v> of type <%T>", v, v)
		}
		mappedTokens[key] = v
	}

	s.storage.tokens = mappedTokens

	return nil
}

// GetAppToken gets a token used by an outside service
func (s *SecretManager) GetAppToken(app string) (string, error) {
	data := s.storage.settings[app]

	if data == "" {
		return "", errors.New("app " + app + " not found in SecretManager")
	}

	return data, nil
}

// SyncAppSettings syncs all application settings from Vault and saves them locally
func (s *SecretManager) SyncAppSettings() error {
	data, err := s.fetchData("settings")

	if err != nil {
		log.Println(err)
		return err
	}

	mappedSettings := make(map[string]string, len(data))

	for key, value := range data {
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("invalid conversion to string for value <%v> of type <%T>", v, v)
		}
		mappedSettings[key] = v
	}

	s.storage.settings = mappedSettings

	return nil
}

// GetSetting returns a single setting value
func (s *SecretManager) GetSetting(key string) (string, error) {
	data := s.storage.settings[key]

	if data == "" {
		return "", errors.New("key " + key + " not found in SecretManager")
	}

	return data, nil
}

func (s *SecretManager) fetchData(subpath string) (map[string]interface{}, error) {
	response, err := s.client.Logical().Read(namespace + "/" + subpath)

	if err != nil {
		return nil, err
	}

	if response == nil {
		return nil, errors.New("empty response from Vault")
	}

	return response.Data, nil
}
