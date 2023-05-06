package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/jandro-es/merlin/configs"
	"github.com/jandro-es/merlin/models"
)

func ConfigurableHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Get the endpoint configuration based on the request path and method
		endpointConfig, ok := configs.FindConfiguration(r.Method, r.URL.Path)
		if !ok {
			http.Error(rw, "Endpoint not found", http.StatusNotFound)
			return
		}

		// Gets the params of the request
		params := r.URL.Query()

		// Validates that all the parameters are present
		// TODO: Check for type
		for _, p := range endpointConfig.Parameters {
			if params.Get(p) == "" {
				http.Error(rw, fmt.Sprintf("Missing required parameter '%s'", p), http.StatusBadRequest)
				return
			}
		}
		// Results map from any subrequest
		results := make(map[string]map[string]interface{})

		subRequests := r.Context().Value("subRequests").([]string)
		for _, value := range subRequests {
			results[value] = r.Context().Value(value).(map[string]interface{})
		}

		// Convert response payload to JSON
		responsePayloadJSON, err := json.Marshal(generateResponse(endpointConfig, params, results))
		if err != nil {
			http.Error(rw, "Failed to marshal response payload to JSON", http.StatusInternalServerError)
			return
		}
		rw.Write(responsePayloadJSON)
	}
}

func generateResponse(requestConfig models.EndpointConfig, params url.Values, results map[string]map[string]interface{}) map[string]interface{} {
	// Build response payload
	responsePayload := make(map[string]interface{})
	for key, config := range requestConfig.Response.Values {
		if config.Passthrough {
			// Value is passed directly for the original request
			responsePayload[key] = params.Get(key)
		} else {
			switch config.Generation.Type {
			case "subrequest":
				responsePayload[key] = results[config.Generation.Origin][config.Generation.Field]
			default:
				log.Fatalf("The generation type for the value %s is not supported: %s", key, config.Generation.Type)
				os.Exit(1)
			}
		}
	}
	return responsePayload
}
