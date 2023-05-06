package models

type EndpointConfig struct {
	Path         string                      `yaml:"path"`
	Method       string                      `yaml:"method"`
	Auth         bool                        `yaml:"auth"`
	AuthProvider string                      `yaml:"auth_provider,omitempty"`
	Headers      map[string]RequestHeader    `yaml:"headers"`
	SubRequests  map[string]SubRequestConfig `yaml:"requests,omitempty"`
	Parameters   []string                    `yaml:"parameters,omitempty"`
	Response     ResponseConfig              `yaml:"response"`
}

type RequestHeader struct {
	Required    bool                    `yaml:"required"`
	Passthrough bool                    `yaml:"passthrough"`
	Validation  RequestHeaderValidation `yaml:"validation"`
}

type RequestHeaderValidation struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value,omitempty"`
}

type EndpointConfigs struct {
	Endpoints map[string]EndpointConfig
}
