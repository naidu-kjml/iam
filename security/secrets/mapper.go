package secrets

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"gitlab.skypicker.com/platform/security/iam/shared"
	"gopkg.in/yaml.v3"
)

// Mapper holds the current service secrets configuration
type Mapper struct {
	config Configuration
}

const configurationFile string = "config/secrets/secrets.yaml"

// ServiceConfiguration represents a single service's secret mapping
type ServiceConfiguration struct {
	Name                string   `yaml:"serviceName"`
	AllowedEnvironments []string `yaml:"allowedEnvironments"`
	KeyBase             string   `yaml:"KeyBase"`
}

// Configuration represents services' secrets configuration
type Configuration struct {
	Mappigns []ServiceConfiguration `yaml:"mappings"`
}

// CreateNewConfigurationMapper initializes a mapper from a local yaml configuration
func CreateNewConfigurationMapper() (*Mapper, error) {
	config, err := readConfiguration()
	if err != nil {
		return nil, err
	}

	return &Mapper{config}, nil
}

// GetKeyName returns a service's key name based on the current mapping
func (m Mapper) GetKeyName(serviceName, environment string) (string, error) {
	for _, service := range m.config.Mappigns {
		if strings.EqualFold(serviceName, service.Name) {
			if environment == "" {
				return service.KeyBase, nil
			}
			if shared.StringInSlice(environment, service.AllowedEnvironments) {
				return service.KeyBase + "_" + strings.ToUpper(environment), nil
			}

			return "", errors.New("environment '" + environment + "' is not allowed")
		}
	}

	// If there are no mappings use the service name in upper case
	return strings.ToUpper(serviceName), nil
}

func readConfiguration() (Configuration, error) {
	file, err := os.Open(configurationFile)
	if err != nil {
		return Configuration{}, errors.Wrap(err, "error opening configuration file")
	}

	var configuration Configuration
	yd := yaml.NewDecoder(file)
	if err = yd.Decode(&configuration); err != nil {
		return Configuration{}, errors.Wrap(err, "error reading configuration file")
	}

	return configuration, nil
}
