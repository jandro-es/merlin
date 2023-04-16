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
		for _, p := range endpointConfig.Params {
			if params.Get(p) == "" {
				http.Error(rw, fmt.Sprintf("Missing required parameter '%s'", p), http.StatusBadRequest)
				return
			}
		}

		// Generate the response JSON
		response := make(map[string]interface{})
		for _, field := range endpointConfig.Response.Fields {
			response[field] = params.Get(field)
		}
		jsonBytes, err := json.Marshal(response)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// Send the response
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(jsonBytes)
	}
}
