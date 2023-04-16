package models

type EndpointConfig struct {
	Path     string
	Method   string
	Params   []string
	Auth     string
	Response struct {
		Fields []string
	}
}

type EndpointConfigs struct {
	Endpoints map[string]EndpointConfig
}
