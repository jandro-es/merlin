package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jandro-es/merlin/configs"
	"github.com/jandro-es/merlin/helpers"
	"github.com/jandro-es/merlin/models"
)

func Subrequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !strings.Contains(r.URL.Path, "/api/") {
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		// Get the endpoint configuration based on the request path and method
		endpointConfig, ok := configs.FindConfiguration(r.Method, r.URL.Path)
		if !ok {
			http.Error(w, "Endpoint not found", http.StatusNotFound)
			return
		}

		// We get all the different subrequests for the endpoint
		requests := endpointConfig.SubRequests
		if len(requests) != 0 {
			// Create a wait group object
			var wg sync.WaitGroup
			// Create a mutex object to protect the results slice
			var mu sync.Mutex
			// Create a slice to hold the response objects
			results := make(map[string]map[string]interface{})

			buildAndExecuteRequests := func(requests map[string]models.SubRequestConfig) {
				// We loop through the SubRequestConfigs
				for key, requestConfig := range requests {
					values := parseParameterValues(requestConfig, r)
					wg.Add(1)
					go executeRequest(key, requestConfig, values, r, &wg, &mu, &results)
				}
			}
			buildAndExecuteRequests(requests)
			// Wait for all the HTTP requests to complete
			wg.Wait()

			keys := make([]string, 0, len(results))
			for key, result := range results {
				keys = append(keys, key)
				ctx = context.WithValue(ctx, key, result)
			}
			ctx = context.WithValue(ctx, "subRequests", keys)

			// Pass the context object to the next middleware and the handler
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			// No subrequests to perform
			ctx := r.Context()
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}

func parseParameterValues(requestConfig models.SubRequestConfig, r *http.Request) map[string]string {
	values := make(map[string]string)
	for key, parameterConfig := range requestConfig.Parameters {
		var variable = extractVariable(parameterConfig.Value)
		if variable != nil {
			values[key] = processVariable(*variable, r)
		} else {
			// No variable found to replace, we add the original value.
			values[key] = parameterConfig.Value
		}
	}
	return values
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

func processVariable(variable string, r *http.Request) string {
	// We extract the indicator for the origin of the value.
	parts := strings.Split(variable, "_")
	if len(parts) < 2 {
		log.Fatalf("The variable %s can't be parsed to the right format", variable)
	}
	origin := parts[0]
	identifier := parts[1]
	var value string
	switch origin {
	case "src":
		// The source of the value is matched on the original request
		value = r.URL.Query().Get(identifier)
		if value == "" {
			// Empty string, the system should fail
			log.Fatalf("No value was found for the identifier: %s", identifier)
			os.Exit(1)
		}
	default:
		log.Fatalf("The origin of the value is not supported: %s", origin)
		os.Exit(1)
	}
	return value
}

func executeRequest(key string, requestConfig models.SubRequestConfig, values map[string]string, r *http.Request, wg *sync.WaitGroup, mutex *sync.Mutex, results *map[string]map[string]interface{}) {
	// We need to replace the URL parameters if there are any
	re := regexp.MustCompile("<([^>]+)>")
	parsedUrl := re.ReplaceAllStringFunc(requestConfig.Path, func(match string) string {
		key := match[1 : len(match)-1]
		value, ok := values[key]
		if ok {
			encodedValue := url.QueryEscape(value)
			return encodedValue
		}
		return match
	})
	// TODO: do the same for query string parameters

	defer wg.Done()

	client := &http.Client{
		Timeout: time.Second * 15,
	}
	req, err := http.NewRequest(requestConfig.Method, parsedUrl, nil)
	if err != nil {
		helpers.ExitOnFail(err, "Error while creating the request for the subrequest")
	}
	// We need to add the relevant headers as per the definition
	setSubRequestHeaders(req, requestConfig, r)
	// Send the HTTP request
	resp, err := client.Do(req)

	// TODO: Manage error properly
	if err != nil {
		fmt.Printf("Error fetching %s: %s\n", parsedUrl, err.Error())
		return
	}

	// Lock the results slice to add the response body
	mutex.Lock()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Print("Error parsing the response")
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		// TODO: Proper error handling
	}
	(*results)[key] = result
	mutex.Unlock()
}

func setSubRequestHeaders(req *http.Request, requestConfig models.SubRequestConfig, r *http.Request) {
	req.Header.Set("user-agent", "merlin")
	req.Header.Set("Content-Type", "application/json")
	for key, headerConfig := range requestConfig.Headers {
		if headerConfig.Passthrough {
			req.Header.Set(key, r.Header.Get(key))
		} else {
			switch headerConfig.Generation.Type {
			case "uuid":
				req.Header.Set(key, uuid.New().String())
			default:
				log.Fatalf("The generation type for the header %s is not supported: %s", key, headerConfig.Generation.Type)
				os.Exit(1)
			}
		}
	}

	// Special case is the authorization header, if auth is required we need to passthough the athorization
	// header from the request.
	if requestConfig.Auth {
		req.Header.Set("Authorization", r.Header.Get("Authorization"))
	}
}
