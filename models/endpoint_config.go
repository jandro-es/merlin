package models

type EndpointConfig struct {
	Path       string                   `yaml:"path"`
	Method     string                   `yaml:"method"`
	Headers    map[string]RequestHeader `yaml:"headers"`
	Parameters []string                 `yaml:"parameters,omitempty"`
	Auth       bool                     `yaml:"auth,omitempty"`
	Response   ResponseConfig           `yaml:"response"`
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

type ResponseConfig struct {
	Headers map[string]string `yaml:"headers,omitempty"`
	Body    interface{}       `yaml:"body"`
}

type EndpointConfigs struct {
	Endpoints map[string]EndpointConfig
}
