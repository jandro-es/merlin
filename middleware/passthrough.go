package middleware

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/jandro-es/merlin/configs"
)

// Middleware function to pass the origin headers to the response as defined.
func PassthroughHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the endpoint configuration based on the request path and method
		endpointConfig, ok := configs.FindConfiguration(r.Method, r.URL.Path)
		if !ok {
			http.Error(w, "Endpoint not found", http.StatusNotFound)
			return
		}

		for _, value := range endpointConfig.Response.Headers {
			variable := extractVariable(value)
			if variable != nil {
				passthroughKey, passthroughValue := processVariable(*variable, r)
				w.Header().Set(passthroughKey, passthroughValue)
			}
			// If no variable exists, we don't need to add anything
		}

		next.ServeHTTP(w, r)
	})
}

// Function to extract the variable name from a response value. The format is {{<something>}}. If not found
// the system returns nil for the parser to return the direct value.
func extractVariable(key string) *string {
	re := regexp.MustCompile(`{{(.+?)}}`)
	// Find the first match
	match := re.FindStringSubmatch(key)
	// Check if a match was found
	if len(match) > 1 {
		// Extract the captured group and return it
		captured := match[1]
		return &captured
	} else {
		// No variable was found, returning nil
		return nil
	}
}

func processVariable(variable string, r *http.Request) (key string, value string) {
	// We extract the indicator for the origin of the value.
	parts := strings.Split(variable, "_")
	if len(parts) < 2 {
		log.Fatalf("The variable %s can't be parsed to the right format", variable)
	}
	origin := parts[0]
	identifier := parts[1]
	var passthroughValue string
	var passthroughKey string
	switch origin {
	case "src":
		// we need to find the key in the original config for the extracted identifier
		for name, values := range r.Header {
			fmt.Printf("Header key: %s", name)
			fmt.Println()
			for _, value := range values {
				fmt.Printf("Header value: %s", value)
				fmt.Println()
				if extractVariable(value) == &name {
					passthroughValue = r.Header.Get(identifier)
					passthroughKey = name
					fmt.Printf("Header value: %s", passthroughValue)
					fmt.Println()
					fmt.Printf("Header key: %s", passthroughKey)
					fmt.Println()
				}
			}
		}
		if passthroughValue == "" && passthroughKey == "" {
			// Empty string, the system should fail
			log.Fatalf("No value was found for the identifier: %s", identifier)
			os.Exit(1)
		}
	default:
		log.Fatalf("The origin of the value is not supported: %s", origin)
		os.Exit(1)
	}
	return passthroughKey, passthroughValue
}
