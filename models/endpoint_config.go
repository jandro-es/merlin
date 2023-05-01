package models

type EndpointConfig struct {
	Path        string                      `yaml:"path"`
	Method      string                      `yaml:"method"`
	Headers     map[string]RequestHeader    `yaml:"headers"`
	SubRequests map[string]SubRequestConfig `yaml:"requests,omitempty"`
	Parameters  []string                    `yaml:"parameters,omitempty"`
	Auth        bool                        `yaml:"auth,omitempty"`
	Response    ResponseConfig              `yaml:"response"`
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
