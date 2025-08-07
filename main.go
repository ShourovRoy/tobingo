// Package tobingo provides a lightweight HTTP router with path parameter support
package tobingo

import (
	"context"
	"net/http"
	"strings"
)

// Route represents a single HTTP route configuration
type Route struct {
	Method  string             // HTTP method (GET, POST, PUT, DELETE, etc.)
	Path    string             // URL path pattern, can include parameters like "/users/:id"
	Handler http.HandlerFunc   // Handler function to execute when route matches
}

// contextKey is a custom type used for context keys to avoid collisions
type contextKey string

// ParamsKey is the context key used to store path parameters in the request context
const ParamsKey contextKey = "params"

// Rastauter is the main router struct that holds all registered routes
type Rastauter struct {
	routes []Route // Slice containing all registered routes
}

// NewRastaRouterInitializer creates and returns a new instance of Rastauter
// with an empty routes slice ready for route registration
func NewRastaRouterInitializer() *Rastauter {
	return &Rastauter{
		routes: []Route{},
	}
}

// StartServer starts the HTTP server on the specified port using this router
// The port should be in the format ":8080" or "localhost:8080"
// Returns an error if the server fails to start
func (rt *Rastauter) StartServer(port string) error {
	return http.ListenAndServe(port, rt)
}

// GET registers a new GET route with the specified path pattern and handler
// Path can include parameters using colon notation (e.g., "/users/:id")
// The handler will be called when a GET request matches the path pattern
func (rt *Rastauter) GET(path string, handler http.HandlerFunc) {
	rt.routes = append(rt.routes, Route{
		Method: "GET",
		Path: path,
		Handler: handler,
	})
}

// GetParam extracts a path parameter value from the request context
// Returns the parameter value if found, otherwise returns an empty string
// Example: For route "/users/:id" and request "/users/123", GetParam(r, "id") returns "123"
func GetParam(r *http.Request, key string) string {
	// Attempt to retrieve parameters map from request context
	if params, ok := r.Context().Value(ParamsKey).(map[string]string); ok {
		return params[key]
	}
	return ""
}

// ServeHTTP implements the http.Handler interface, making Rastauter compatible with net/http
// This method is called for every HTTP request and handles route matching and parameter extraction
func (rt *Rastauter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	
	// Iterate through all registered routes to find a match
	for _, route := range rt.routes {
		// First check if the HTTP method matches
		if route.Method == r.Method {

			// Get the route path pattern from the registered route
			routePath := route.Path
			
			// Split the route path into segments, trimming spaces and splitting by "/"
			// Note: This trims spaces instead of "/" which might be intentional
			routerPathSlice := strings.Split(strings.Trim(routePath, " "), "/")
			
			// Get the actual request path, trim trailing "/" and split into segments
			requestPath := strings.Trim(r.URL.Path, "/")
			requestPathSlice := strings.Split(requestPath, "/")

			// Check if the number of path segments match
			// Skip the first element in routerPathSlice with [1:] (assumes it's empty from leading "/")
			if len(routerPathSlice[1:]) != len(requestPathSlice) {
				continue // Try next route if segment count doesn't match
			}
			
			// Initialize map to store extracted path parameters
			params := make(map[string]string)
			
			// Iterate through each segment of the route pattern
			for routerIndex, routerPathName := range routerPathSlice[1:] {
				// Check if this segment is a parameter (starts with ":")
				if after, ok := strings.CutPrefix(routerPathName, ":"); ok {
					// Extract parameter name (everything after ":")
					paramName := after
					// Store the corresponding value from the request path
					params[paramName] = requestPathSlice[routerIndex]
				}
				// Note: This implementation doesn't validate exact matches for non-parameter segments
				// All routes with matching segment counts will match, regardless of literal segment values
			}

			// Add the extracted parameters to the request context
			// This makes them available to the handler via GetParam function
			ctx := context.WithValue(r.Context(), ParamsKey, params)
			r = r.WithContext(ctx)
			
			// Execute the matched route's handler
			route.Handler(w, r)

			// Return early since we found a match and handled the request
			return
		}
	}
	
	// If no route matches the request method and path, return 404 Not Found
	http.NotFound(w, r)
}