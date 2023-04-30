package configs

import (
	"io/ioutil"
	"log"

	"github.com/jandro-es/merlin/models"
	"gopkg.in/yaml.v3"
)

var Configurations *models.EndpointConfigs = ParseConfigurations()

func ParseConfigurations() *models.EndpointConfigs {
	const path = "./endpoints"
	var endpointConfigs models.EndpointConfigs
	endpointConfigs.Endpoints = make(map[string]models.EndpointConfig)

	configFiles, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatalf("Failed while loading endpoint configurations with error %s", err)
	}

	for _, file := range configFiles {
		configData, err := ioutil.ReadFile(path + "/" + file.Name())
		if err != nil {
			log.Fatalf("Failed to read the file %s with error %s", file.Name(), err)
		}
		config, err := parseEndpointConfigurations(configData)

		if err != nil {
			log.Fatalf("Failed to parse the configuration for %s with error %s", file.Name(), err)
		}
		// Add the endpoint config to the map
		key := config.Method + config.Path
		endpointConfigs.Endpoints[key] = config
	}
	return &endpointConfigs
}

// Get the endpoint configuration based on the request path and method
func FindConfiguration(method string, path string) (models.EndpointConfig, bool) {
	endpointConfig, ok := Configurations.Endpoints[method+path]
	if !ok {
		return models.EndpointConfig{}, false
	}
	return endpointConfig, true
}

func parseEndpointConfigurations(data []byte) (models.EndpointConfig, error) {
	var config models.EndpointConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return models.EndpointConfig{}, err
	}
	return config, nil
}
