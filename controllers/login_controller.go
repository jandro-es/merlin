package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jandro-es/merlin/configs"
	"github.com/jandro-es/merlin/models"
)

func Login() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		var creds models.Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
			return
		}

		firebaseAuth := configs.FirebaseClient

		token, err := firebaseAuth.CustomToken(context.Background(), "Firebase UUID")
		if err != nil {
			http.Error(rw, fmt.Sprintf("unable to generate token with error: %s", err), http.StatusUnauthorized)
			panic(err)
		}

		// Print the Firebase token
		fmt.Println("Firebase token:", token)

		// Marshal the token to JSON
		tokenJSON, err := json.Marshal(token)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the token JSON to the response body
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		rw.Write(tokenJSON)
	}
}
