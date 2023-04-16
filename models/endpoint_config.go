package models

type EndpointConfig struct {
	Path       string         `yaml:"path"`
	Method     string         `yaml:"method"`
	Parameters []string       `yaml:"parameters,omitempty"`
	Auth       bool           `yaml:"auth,omitempty"`
	Response   ResponseConfig `yaml:"response"`
}

type ResponseConfig struct {
	Headers map[string]string `yaml:"headers,omitempty"`
	Body    interface{}       `yaml:"body"`
}

type EndpointConfigs struct {
	Endpoints map[string]EndpointConfig
}
