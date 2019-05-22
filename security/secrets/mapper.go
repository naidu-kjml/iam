package secrets

import (
	"os"
	"strings"

	"github.com/pkg/errors"
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
	TokenBase           string   `yaml:"tokenBase"`
}

// Configuration represents services' secrets configuration
type Configuration struct {
	Mappings []ServiceConfiguration `yaml:"mappings"`
}

// CreateNewConfigurationMapper initializes a mapper from a local yaml configuration
func CreateNewConfigurationMapper() (*Mapper, error) {
	config, err := readConfiguration()
	if err != nil {
		return nil, err
	}

	return &Mapper{config}, nil
}

// GetTokenName returns a service's token name based on the current mapping
func (m Mapper) GetTokenName(serviceName, environment string) (string, error) {
	for _, service := range m.config.Mappings {
		if strings.EqualFold(serviceName, service.Name) {
			if environment == "" {
				return service.TokenBase, nil
			}
			if stringInSlice(environment, service.AllowedEnvironments) {
				return service.TokenBase + "_" + strings.ToUpper(environment), nil
			}

			return "", errors.New("environment '" + environment + "' is not allowed")
		}
	}

	return "", errors.New("unknown service '" + serviceName + "'")
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
