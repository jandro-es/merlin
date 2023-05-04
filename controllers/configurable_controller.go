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

		// // Generate the response JSON
		// response := make(map[string]interface{})
		// for _, field := range endpointConfig.Response.Body {
		// 	response[field] = params.Get(field)
		// }
		// jsonBytes, err := json.Marshal(response)
		// if err != nil {
		// 	http.Error(rw, err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		// Set the headers as per the configuration
		// for key, value := range endpointConfig.Response.Headers {
		// 	rw.Header().Set(key, value)
		// }
		// Send the response
		// rw.Write(jsonBytes)
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

// func parseValue(value interface{}, params url.Values) interface{} {
// 	switch v := value.(type) {
// 	case string:
// 		var variable = extractVariable(v)
// 		if variable != nil {
// 			return processVariable(*variable, params)
// 		} else {
// 			// No variable found to replace, returning original value
// 			return v
// 		}
// 	default:
// 		return value
// 	}
// }

// // Function to extract the variable name from a response value. The format is {{<something>}}. If not found
// // the system returns nil for the parser to return the direct value.
// func extractVariable(key string) *string {
// 	re := regexp.MustCompile(`{{(.+?)}}`)
// 	// Find the first match
// 	match := re.FindStringSubmatch(key)
// 	// Check if a match was found
// 	if len(match) > 1 {
// 		// Extract the captured group and return it
// 		captured := match[1]
// 		return &captured
// 	} else {
// 		// No variable was found, returning nil
// 		return nil
// 	}
// }

// // TODO: instead of failing, it should return an error status and then the handler return an specific HTTP code.
// func processVariable(variable string, params url.Values) interface{} {
// 	// We extract the indicator for the origin of the value.
// 	parts := strings.Split(variable, "_")
// 	if len(parts) < 2 {
// 		log.Fatalf("The variable %s can't be parsed to the right format", variable)
// 	}
// 	origin := parts[0]
// 	identifier := parts[1]
// 	var value string
// 	switch origin {
// 	case "src":
// 		// The source of the value is mathed on the original request
// 		value = params.Get(identifier)
// 		if value == "" {
// 			// Empty string, the system should fail
// 			log.Fatalf("No value was found for the identifier: %s", identifier)
// 			os.Exit(1)
// 		}
// 	default:
// 		log.Fatalf("The origin of the value is not supported: %s", origin)
// 		os.Exit(1)
// 	}
// 	return value
// }
