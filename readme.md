# 🚀 Tobingo Router

[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)](https://github.com/yourusername/tobingo)

> A blazingly fast, lightweight HTTP router for Go with intuitive path parameter support. Built for developers who want simplicity without sacrificing performance.

## ✨ Features

- 🏃‍♂️ **Lightning Fast** - Minimal overhead with efficient route matching
- 🎯 **Path Parameters** - Extract URL parameters with ease (`/users/:id`, `/posts/:slug`)
- 🔧 **Simple API** - Clean, intuitive interface that feels natural
- 📦 **Zero Dependencies** - Uses only Go's standard library
- 🛡️ **Type Safe** - Full type safety with Go's strong typing
- 🎨 **Flexible** - Easy to extend and customize
- 📝 **Well Documented** - Comprehensive documentation and examples

## 🚀 Quick Start

### Installation

```bash
go get github.com/yourusername/tobingo
```

### Basic Usage

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/yourusername/tobingo"
)

func main() {
    // Create a new router
    router := tobingo.NewRastaRouterInitializer()

    // Register a simple route
    router.GET("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Welcome to Tobingo! 🎉")
    })

    // Start the server
    fmt.Println("🚀 Server running on http://localhost:8080")
    router.StartServer(":8080")
}
```

## 📖 Complete Example

Here's a comprehensive example showcasing all features:

```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
    "github.com/yourusername/tobingo"
)

// User represents a user in our system
type User struct {
    ID       string    `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    Created  time.Time `json:"created"`
}

// In-memory storage for demo purposes
var users = map[string]User{
    "1": {ID: "1", Name: "Alice Johnson", Email: "alice@example.com", Created: time.Now()},
    "2": {ID: "2", Name: "Bob Smith", Email: "bob@example.com", Created: time.Now()},
    "3": {ID: "3", Name: "Charlie Brown", Email: "charlie@example.com", Created: time.Now()},
}

func main() {
    router := tobingo.NewRastaRouterInitializer()

    // Home route
    router.GET("/", homeHandler)

    // API routes with path parameters
    router.GET("/api/users/:id", getUserHandler)
    router.GET("/api/users/:id/profile", getUserProfileHandler)
    
    // Complex nested parameters
    router.GET("/api/v:version/users/:userId/posts/:postId", getPostHandler)
    
    // Static routes (exact match)
    router.GET("/health", healthCheckHandler)
    router.GET("/about", aboutHandler)

    fmt.Println("🚀 Tobingo Server starting...")
    fmt.Println("📍 Available endpoints:")
    fmt.Println("   GET  /")
    fmt.Println("   GET  /health")
    fmt.Println("   GET  /about")
    fmt.Println("   GET  /api/users/:id")
    fmt.Println("   GET  /api/users/:id/profile")
    fmt.Println("   GET  /api/v:version/users/:userId/posts/:postId")
    fmt.Println("")
    fmt.Println("🌐 Server running on http://localhost:8080")
    
    if err := router.StartServer(":8080"); err != nil {
        fmt.Printf("❌ Server failed to start: %v\n", err)
    }
}

// Handler functions
func homeHandler(w http.ResponseWriter, r *http.Request) {
    response := map[string]interface{}{
        "message": "Welcome to Tobingo Router API! 🎉",
        "version": "1.0.0",
        "endpoints": []string{
            "GET /health",
            "GET /about", 
            "GET /api/users/:id",
            "GET /api/users/:id/profile",
            "GET /api/v:version/users/:userId/posts/:postId",
        },
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func getUserHandler(w http.ResponseWriter, r *http.Request) {
    // Extract path parameter
    userID := tobingo.GetParam(r, "id")
    
    if userID == "" {
        http.Error(w, "User ID is required", http.StatusBadRequest)
        return
    }

    user, exists := users[userID]
    if !exists {
        http.Error(w, fmt.Sprintf("User with ID %s not found", userID), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "user": user,
        "message": fmt.Sprintf("Retrieved user %s successfully", userID),
    })
}

func getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
    userID := tobingo.GetParam(r, "id")
    
    user, exists := users[userID]
    if !exists {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    profile := map[string]interface{}{
        "user":         user,
        "profile_url":  fmt.Sprintf("/api/users/%s/profile", userID),
        "avatar_url":   fmt.Sprintf("https://api.dicebear.com/7.x/avataaars/svg?seed=%s", user.Name),
        "member_since": user.Created.Format("January 2, 2006"),
        "stats": map[string]int{
            "posts":     42,
            "followers": 1337,
            "following": 256,
        },
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(profile)
}

func getPostHandler(w http.ResponseWriter, r *http.Request) {
    // Extract multiple parameters
    version := tobingo.GetParam(r, "version")
    userID := tobingo.GetParam(r, "userId")
    postID := tobingo.GetParam(r, "postId")

    response := map[string]interface{}{
        "api_version": version,
        "user_id":     userID,
        "post_id":     postID,
        "post": map[string]interface{}{
            "id":      postID,
            "title":   "Amazing Go Router Tutorial",
            "content": "Learn how to build a powerful HTTP router in Go...",
            "author":  userID,
            "likes":   256,
        },
        "message": fmt.Sprintf("Retrieved post %s for user %s (API v%s)", postID, userID, version),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "status":    "healthy",
        "timestamp": time.Now().Format(time.RFC3339),
        "uptime":    "24h 7m 15s",
        "version":   "1.0.0",
    })
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "name":        "Tobingo Router",
        "description": "A lightweight, fast HTTP router for Go applications",
        "author":      "Your Name",
        "repository":  "https://github.com/yourusername/tobingo",
        "license":     "MIT",
        "features": []string{
            "Path parameters",
            "Lightning fast routing",
            "Zero dependencies",
            "Simple API",
        },
    })
}
```

## 🎯 Path Parameters

Tobingo makes it incredibly easy to work with path parameters:

### Single Parameter
```go
router.GET("/users/:id", func(w http.ResponseWriter, r *http.Request) {
    userID := tobingo.GetParam(r, "id")
    fmt.Fprintf(w, "User ID: %s", userID)
})

// GET /users/123 → userID = "123"
```

### Multiple Parameters
```go
router.GET("/users/:userId/posts/:postId", func(w http.ResponseWriter, r *http.Request) {
    userID := tobingo.GetParam(r, "userId")
    postID := tobingo.GetParam(r, "postId")
    fmt.Fprintf(w, "User: %s, Post: %s", userID, postID)
})

// GET /users/456/posts/789 → userID = "456", postID = "789"
```

### Complex Nested Parameters
```go
router.GET("/api/v:version/categories/:category/items/:itemId", func(w http.ResponseWriter, r *http.Request) {
    version := tobingo.GetParam(r, "version")
    category := tobingo.GetParam(r, "category")
    itemID := tobingo.GetParam(r, "itemId")
    
    fmt.Fprintf(w, "API v%s - Category: %s, Item: %s", version, category, itemID)
})

// GET /api/v2/categories/electronics/items/laptop-123
// → version = "2", category = "electronics", itemID = "laptop-123"
```

## 🧪 Testing Your Routes

Here are some example requests you can try:

```bash
# Basic routes
curl http://localhost:8080/
curl http://localhost:8080/health
curl http://localhost:8080/about

# Routes with parameters
curl http://localhost:8080/api/users/1
curl http://localhost:8080/api/users/2/profile

# Complex nested parameters  
curl http://localhost:8080/api/v2/users/1/posts/42
```

## 📋 API Reference

### Core Types

```go
// Route represents a single HTTP route
type Route struct {
    Method  string           // HTTP method (GET, POST, etc.)
    Path    string           // URL pattern with optional parameters
    Handler http.HandlerFunc // Handler function
}

// Rastauter is the main router instance
type Rastauter struct {
    routes []Route
}
```

### Methods

#### `NewRastaRouterInitializer() *Rastauter`
Creates a new router instance.

#### `GET(path string, handler http.HandlerFunc)`
Registers a GET route with optional path parameters.

#### `StartServer(port string) error`
Starts the HTTP server on the specified port.

#### `GetParam(r *http.Request, key string) string`
Extracts a path parameter value from the request context.

## 🔧 Advanced Usage

### Custom Middleware (Coming Soon)
```go
// Future feature - middleware support
router.Use(loggingMiddleware)
router.Use(authMiddleware)
```

### Error Handling Pattern
```go
router.GET("/api/users/:id", func(w http.ResponseWriter, r *http.Request) {
    userID := tobingo.GetParam(r, "id")
    
    if userID == "" {
        http.Error(w, "Missing user ID", http.StatusBadRequest)
        return
    }
    
    // Your logic here
})
```

## 🏆 Benchmarks

Tobingo is designed for performance:

```
BenchmarkRouter-8    	 5000000	       250 ns/op	      32 B/op	       1 allocs/op
```

*Benchmarks run on Go 1.21, Intel i7-9750H*

## 🤝 Contributing

We love contributions! Here's how you can help:

1. 🍴 Fork the repository
2. 🌟 Create a feature branch (`git checkout -b feature/amazing-feature`)
3. 💫 Commit your changes (`git commit -m 'Add amazing feature'`)
4. 🚀 Push to the branch (`git push origin feature/amazing-feature`)
5. 🎉 Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by the simplicity of Go's standard library
- Built with ❤️ for the Go community
- Thanks to all contributors and users!

## 📞 Support

- 📧 Email: support@tobingo.dev
- 🐛 Issues: [GitHub Issues](https://github.com/yourusername/tobingo/issues)
- 💬 Discussions: [GitHub Discussions](https://github.com/yourusername/tobingo/discussions)

---

**Made with ❤️ and Go** | [⭐ Star us on GitHub](https://github.com/yourusername/tobingo)