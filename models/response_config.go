package models

type ResponseConfig struct {
	Headers map[string]ResponseHeader `yaml:"headers"`
	Values  map[string]ResponseValue  `yaml:"values"`
}

type ResponseHeader struct {
	Passthrough bool               `yaml:"passthrough"`
	Generation  ResponseGeneration `yaml:"generation"`
}

type ResponseGeneration struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value,omitempty"`
}

type ResponseValue struct {
	Passthrough bool                    `yaml:"passthrough"`
	Generation  ResponseValueGeneration `yaml:"generation"`
}

type ResponseValueGeneration struct {
	Type   string `yaml:"type"`
	Origin string `yaml:"origin,omitempty"`
	Field  string `yaml:"field,omitempty"`
}
