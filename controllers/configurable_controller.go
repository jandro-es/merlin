package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jandro-es/merlin/configs"
)

func ConfigurableHandler() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		// Get the endpoint configuration based on the request path and method
		endpointConfig, ok := configs.Configurations.Endpoints[r.Method+r.URL.Path]
		if !ok {
			http.Error(rw, "Endpoint not found", http.StatusNotFound)
			return
		}

		// Get the request parameters
		// params := mux.Vars(r)
		params := r.URL.Query()

		// Debug only
		// bs, _ := json.Marshal(params)
		// fmt.Println(string(bs))

		// Check for required parameters
		for _, p := range endpointConfig.Parameters {
			if params.Get(p) == "" {
				http.Error(rw, fmt.Sprintf("Missing required parameter '%s'", p), http.StatusBadRequest)
				return
			}
		}

		// Build response payload
		responsePayload := map[string]interface{}{}
		// bs, _ := json.Marshal(endpointConfig.Response.Body)
		// fmt.Println(string(bs))
		// fmt.Printf("%T", endpointConfig.Response.Body)
		switch v := endpointConfig.Response.Body.(type) {
		case map[string]interface{}:
			fmt.Println("CASE 1")
			// bs, _ := json.Marshal(params)
			// fmt.Println(string(bs))
			for key, value := range v {
				responsePayload[key] = value
			}
		case []interface{}:
			fmt.Println("CASE 2")
			responsePayloadSlice := make([]interface{}, len(v))
			for i, value := range v {
				responsePayloadSlice[i] = value
			}
			responsePayload["data"] = responsePayloadSlice
		default:
			fmt.Println("CASE DEFAULT")
			responsePayload = nil
		}

		// Convert response payload to JSON
		responsePayloadJSON, err := json.Marshal(responsePayload)
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
		for key, value := range endpointConfig.Response.Headers {
			rw.Header().Set(key, value)
		}
		// Send the response
		// rw.Write(jsonBytes)
		rw.Write(responsePayloadJSON)
	}
}
