package routes

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jandro-es/merlin/configs"
	"github.com/jandro-es/merlin/controllers"
)

func ConfigurableRoutes(router *mux.Router) {
	for _, endpointConfig := range configs.Configurations.Endpoints {
		router.Handle(endpointConfig.Path, handlers.LoggingHandler(os.Stdout, http.HandlerFunc(controllers.ConfigurableHandler()))).Methods(endpointConfig.Method)
		log.Printf("Route loaded for path %s with HTTP method: %s", endpointConfig.Path, endpointConfig.Method)
	}
}
