package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jandro-es/merlin/responses"
)

func Liveness() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		response := responses.GovernanceResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": "Ping Liveness"}}
		json.NewEncoder(rw).Encode(response)
	}
}

func Readiness() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.WriteHeader(http.StatusOK)
		response := responses.GovernanceResponse{Status: http.StatusOK, Message: "Success", Data: map[string]interface{}{"data": "Ping Readiness"}}
		json.NewEncoder(rw).Encode(response)
	}
}
