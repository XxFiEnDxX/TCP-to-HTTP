# TCP-to-HTTP Server Framework

Welcome to TCP-to-HTTP, a custom HTTP/1.1 server implementation built directly on raw TCP sockets in Go. Think of it as rolling your own HTTP server from scratchâ€”no `net/http` package shortcuts here. This project gives you complete control over every byte that flows through your server, making it perfect for learning HTTP internals or building specialized server applications.

Unlike traditional Go HTTP servers that rely on the standard library's abstractions, this implementation parses HTTP requests byte-by-byte from TCP streams and handles the entire HTTP/1.1 protocol lifecycle manually. It's like building a car engine instead of just turning the key.

## What Makes This Special

This repository serves dual purposes that make it both educational and practical:

**ðŸ”§ Reusable Framework**: The modular design in `internal/` packages provides building blocks for creating custom HTTP-like servers with full protocol control.

**ðŸš€ Production-Ready Server**: The `cmd/httpServer` application demonstrates real-world usage with routing, static file serving, error handling, and HTTP proxying with chunked transfer encoding.

## Project Requirements

Before diving in, make sure your development environment meets these requirements:

- **Go 1.19+**: The project uses modern Go features and standard library APIs
- **Unix-like OS**: Developed and tested on Linux/macOS (Windows should work but isn't explicitly tested)
- **Network Access**: Required for the HTTP proxying features to external services

The beauty of this project lies in its minimal external dependenciesâ€”most functionality is built using Go's standard library.

## Dependencies

This project embraces simplicity with a carefully curated set of dependencies:

### Core Framework Dependencies
- **`net`**: TCP listener and connection management
- **`io`**: Reader/Writer interfaces for stream operations  
- **`bufio`**: Buffered reading with Scanner for efficient line-by-line HTTP parsing

### Application-Specific Dependencies
- **`net/http`**: HTTP client for the `/httpbin/*` proxy endpoints
- **`crypto/sha256`**: SHA256 hash computation for response integrity headers
- **`os`**: File system operations for static asset serving

### Development Dependencies
- **`github.com/stretchr/testify`**: Assertion library for comprehensive unit testing

The minimal dependency footprint means faster builds, fewer security concerns, and easier maintenance.

## Getting Started

Ready to explore HTTP at the protocol level? Here's how to get the server running:

### Building the Server

Navigate to your project directory and build the HTTP server application:

```bash
go build -o httpserver ./cmd/httpServer
```

This creates an executable that combines all the framework components into a ready-to-run server.

### Basic Framework Usage

If you want to build your own server using the TCP-to-HTTP framework, here's the essential pattern:

```go
package main

import (
    "github.com/XxFiEnDxX/TCP-to-HTTP/internal/server"
    "github.com/XxFiEnDxX/TCP-to-HTTP/internal/response"
)

func myHandler(w *response.Writer, req *request.Request) {
    w.WriteStatus(200)
    w.WriteHeader("Content-Type", "text/plain")
    w.WriteBody("Hello from my custom server!")
}

func main() {
    server.Serve(8080, myHandler)
    select {} // Keep the program running
}
```

The framework handles all the TCP connection management and HTTP parsingâ€”you just focus on your application logic.

## Running the App

### Starting the Default Server

Launch the pre-configured HTTP server that showcases all framework capabilities:

```bash
./httpserver
```

The server starts immediately and listens on **port 42069**. You'll see output confirming the server is ready to accept connections.

### Testing the Endpoints

The default server provides several endpoints that demonstrate different HTTP features:

**Basic HTML Response**:
```bash
curl http://localhost:42069/
# Returns: 200 OK with HTML success message
```

**Error Handling Examples**:
```bash
curl http://localhost:42069/yourproblem
# Returns: 400 Bad Request with humorous error message

curl http://localhost:42069/myproblem  
# Returns: 500 Internal Server Error with humorous error message
```

**Static File Serving**:
```bash
curl http://localhost:42069/video
# Serves assets/vim.mp4 with proper Content-Type and Content-Length headers
```

**HTTP Proxying with Chunked Transfer**:
```bash
curl http://localhost:42069/httpbin/get
# Proxies to httpbin.org with Transfer-Encoding: chunked
# Includes X-Content-SHA256 and X-Content-Length trailer headers
```

## Code Examples

### Custom Request Handler

Here's how to create a sophisticated handler that demonstrates the framework's capabilities:

```go
func advancedHandler(w *response.Writer, req *request.Request) {
    switch req.Method {
    case "GET":
        if req.Path == "/api/status" {
            w.WriteStatus(200)
            w.WriteHeader("Content-Type", "application/json")
            w.WriteBody(`{"status": "healthy", "version": "1.0.0"}`)
        } else {
            w.WriteStatus(404)
            w.WriteHeader("Content-Type", "text/plain")
            w.WriteBody("Endpoint not found")
        }
    case "POST":
        // Handle POST requests with request body parsing
        body := req.Body
        w.WriteStatus(201)
        w.WriteHeader("Content-Type", "text/plain")
        w.WriteBody("Created resource with body: " + body)
    default:
        w.WriteStatus(405)
        w.WriteHeader("Allow", "GET, POST")
        w.WriteBody("Method not allowed")
    }
}
```

### Request Parsing Deep Dive

The framework's request parser handles the complexity of HTTP protocol parsing:

```go
// The parser converts raw TCP bytes into structured Request objects
type Request struct {
    Method  string
    Path    string  
    Headers map[string]string
    Body    string
}

// Example of what gets parsed from TCP stream:
// "GET /api/users HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.68.0\r\n\r\n"
// Becomes:
// Request{
//     Method: "GET",
//     Path: "/api/users", 
//     Headers: {"Host": "localhost:42069", "User-Agent": "curl/7.68.0"}
// }
```

### Response Writer Usage

The response writer provides a clean API for HTTP response generation:

```go
func jsonHandler(w *response.Writer, req *request.Request) {
    data := map[string]interface{}{
        "message": "Hello World",
        "timestamp": time.Now().Unix(),
        "client_ip": req.RemoteAddr,
    }
    
    jsonBytes, _ := json.Marshal(data)
    
    w.WriteStatus(200)
    w.WriteHeader("Content-Type", "application/json")
    w.WriteHeader("Cache-Control", "no-cache")
    w.WriteBody(string(jsonBytes))
}
```

## Architecture Overview

The framework follows a clean separation of concerns across four main components:

### Server Component (`internal/server`)
Manages TCP connection lifecycle with a goroutine-per-connection model. Each incoming connection gets its own goroutine for concurrent request handling without blocking.

### Request Parser (`internal/requests`)  
Implements a state machine that converts raw TCP byte streams into structured `Request` objects. Handles HTTP method parsing, header extraction, and body reading according to HTTP/1.1 specifications.

### Response Writer (`internal/response`)
Provides a fluent API for generating properly formatted HTTP responses. Handles status codes, headers, and body content with automatic Content-Length calculation.

### Headers System (`internal/headers`)
Manages HTTP header parsing and validation with case-insensitive storage and retrieval. Supports all standard HTTP headers plus custom header handling.

## Project Structure

The codebase follows Go's standard project layout for maximum clarity:

```
â”œâ”€â”€ cmd/                    # Executable applications
â”‚   â”œâ”€â”€ httpServer/        # Production HTTP server
â”‚   â””â”€â”€ tcplistener/       # Diagnostic TCP listener
â”œâ”€â”€ internal/              # Private framework packages  
â”‚   â”œâ”€â”€ server/           # TCP connection management
â”‚   â”œâ”€â”€ requests/         # HTTP request parsing
â”‚   â”œâ”€â”€ response/         # HTTP response generation
â”‚   â””â”€â”€ headers/          # Header parsing and storage
â”œâ”€â”€ assets/               # Static files (vim.mp4)
â”œâ”€â”€ tmp/                  # Test fixtures and samples
â””â”€â”€ go.mod               # Go module dependencies
```

This structure makes it easy to understand the separation between reusable framework code (`internal/`) and application-specific implementations (`cmd/`).

## Why Build This?

Traditional HTTP servers abstract away the protocol details, which is great for productivity but not for understanding. This project gives you x-ray vision into HTTP:

- **Educational Value**: See exactly how HTTP requests are parsed byte-by-byte
- **Protocol Control**: Handle edge cases and implement custom HTTP behavior  
- **Performance Insights**: Understand the real cost of HTTP parsing and response generation
- **Foundation Knowledge**: Build the skills to debug network issues and optimize server performance

Whether you're a systems programming enthusiast or need fine-grained control over HTTP handling, this framework provides the foundation without sacrificing Go's simplicity and performance.

## Ready to Dive Deeper?

This TCP-to-HTTP framework opens up a world of possibilities for custom server development. From learning HTTP internals to building specialized network applications, you now have the tools to work at the protocol level.

Start by running the example server, explore the endpoints, and then dive into the source code to see how each component works. The modular design makes it easy to understand each piece independently, then see how they work together to create a complete HTTP server.

**Next Steps**: Try modifying the handler functions, add new endpoints, or use the framework components to build your own custom server application. The code is designed to be readable and extensibleâ€”perfect for experimentation and learning.

Happy coding, and enjoy exploring the foundations of web communication! ðŸš€