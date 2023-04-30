package models

type SubRequestConfig struct {
	Path       string                         `yaml:"url"`
	Method     string                         `yaml:"method"`
	Headers    map[string]SubRequestHeader    `yaml:"headers"`
	Parameters map[string]SubRequestParameter `yaml:"parameters,omitempty"`
	Auth       bool                           `yaml:"auth,omitempty"`
}

type SubRequestHeader struct {
	Passthrough bool                 `yaml:"passthrough"`
	Generation  SubRequestGeneration `yaml:"generation"`
}

type SubRequestGeneration struct {
	Type string `yaml:"type"`
}

type SubRequestParameter struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value,omitempty"`
}
