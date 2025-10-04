# TCP-to-HTTP Framework

Ever wondered how HTTP actually works under the hood? This project strips away the magic of Go's standard `net/http` package and builds an HTTP/1.1 server directly on raw TCP sockets. It's like building a car engine from scratch to understand how it really works.

The TCP-to-HTTP framework serves dual purposes: it's both a reusable library for building custom HTTP-like servers and a production-ready HTTP server that demonstrates real-world usage patterns. Think of it as your gateway to understanding the nuts and bolts of web communication.

## What Makes This Special

While most Go HTTP servers rely on the convenience of `net/http`, this implementation parses HTTP requests byte-by-byte from raw TCP streams. It's educational, powerful, and gives you complete control over request parsing and response generation.

The project follows Go's clean architecture principles with a clear separation between reusable framework components and application-specific logic. This means you can use the internal packages to build your own custom servers or study the `cmd/httpServer` example to see how everything fits together.

## Project Requirements

- **Go 1.16+**: This project uses Go modules and requires a modern Go installation
- **Network access**: The HTTP server needs to bind to port 42069 and make outbound connections for proxy features
- **Unix-like environment**: While Go is cross-platform, the project has been primarily tested on Linux and macOS

## Dependencies

This project deliberately uses minimal external dependencies to showcase raw TCP/HTTP implementation:

### Core Dependencies (Framework)
- `net`: TCP listener and connection management
- `io`: Reader/Writer interfaces for I/O operations  
- `bufio`: Buffered reading with Scanner for line-by-line HTTP parsing
- `strings`: String manipulation for HTTP header processing

### Application Dependencies (HTTP Server)
- `net/http`: HTTP client for proxy functionality only
- `crypto/sha256`: SHA256 hashing for response integrity headers
- `os`: File I/O for static file serving

### Test Dependencies
- `github.com/stretchr/testify`: Assertion library for unit tests

The lightweight dependency footprint means you can easily understand and modify every aspect of the HTTP implementation without wrestling with complex external libraries.

## Getting Started

### Running the HTTP Server

The main HTTP server application runs on port 42069 and provides several interesting endpoints:

```bash
go run cmd/httpServer/main.go
```

Once running, you can test the various endpoints:

```bash
# Basic success response
curl http://localhost:42069/

# Error responses with personality
curl http://localhost:42069/yourproblem
curl http://localhost:42069/myproblem

# Static video file serving
curl -I http://localhost:42069/video

# HTTP proxying with chunked encoding
curl http://localhost:42069/httpbin/get
```

### Running the TCP Listener (Diagnostic Tool)

For debugging and learning purposes, you can also run the raw TCP listener:

```bash
go run cmd/tcplistener/main.go
```

This tool accepts raw TCP connections and parses HTTP requests without generating responses, making it perfect for understanding the request parsing process.

## Code Examples

### Basic TCP-to-HTTP Server

Here's how simple it is to create your own HTTP-like server using the framework:

```go
package main

import (
    "fmt"
    "github.com/yourusername/tcp-to-http/internal/server"
    "github.com/yourusername/tcp-to-http/internal/response"
    "github.com/yourusername/tcp-to-http/internal/requests"
)

func myHandler(w *response.Writer, req *requests.Request) {
    if req.Path == "/" {
        w.WriteStatusLine("200 OK")
        w.WriteHeader("Content-Type", "text/plain")
        w.FinishHeaders()
        w.WriteBody("Hello from raw TCP!")
    } else {
        w.WriteStatusLine("404 Not Found")
        w.FinishHeaders()
    }
}

func main() {
    server.Serve(8080, myHandler)
}
```

### Custom Request Processing

The framework gives you full access to parsed HTTP components:

```go
func advancedHandler(w *response.Writer, req *requests.Request) {
    // Access HTTP method
    fmt.Printf("Method: %s\n", req.Method)
    
    // Inspect headers
    contentType := req.Headers.Get("Content-Type")
    userAgent := req.Headers.Get("User-Agent")
    
    // Process request body
    if req.Body != "" {
        fmt.Printf("Body: %s\n", req.Body)
    }
    
    // Generate custom response
    w.WriteStatusLine("200 OK")
    w.WriteHeader("X-Custom-Header", "processed-by-tcp-framework")
    w.WriteHeader("Content-Type", "application/json")
    w.FinishHeaders()
    w.WriteBody(`{"status": "success", "method": "` + req.Method + `"}`)
}
```

### Error Handling and Logging

The framework provides clean error handling patterns:

```go
func handleConnection(conn net.Conn) {
    defer conn.Close()
    
    req, err := requests.RequestFromReader(conn)
    if err != nil {
        log.Printf("Failed to parse request: %v", err)
        return
    }
    
    w := response.NewWriter(conn)
    handler(w, req)
}
```

## Architecture Overview

The project is structured around clean separation of concerns:

### Framework Components (`internal/`)

- **`server/`**: Manages TCP connections and goroutine lifecycle
- **`requests/`**: Parses raw TCP bytes into structured HTTP requests
- **`response/`**: Formats and writes HTTP responses back to clients  
- **`headers/`**: Handles HTTP header parsing and storage

### Applications (`cmd/`)

- **`httpServer/`**: Production HTTP server with routing, static files, and proxying
- **`tcplistener/`**: Diagnostic tool for request parsing analysis

The goroutine-per-connection model ensures excellent concurrency while keeping the code simple and readable. Each client connection runs in its own goroutine, allowing the server to handle hundreds of concurrent requests efficiently.

## Why Build This?

Understanding HTTP at the TCP level gives you superpowers as a web developer. You'll understand why certain performance optimizations work, how to debug network issues, and gain insight into how web frameworks actually function under the hood.

This project is perfect for developers who want to understand the magic behind `net/http`, students learning network programming, or anyone building custom protocols that need HTTP-like request/response patterns.

Ready to dive deep into the world of HTTP? Clone this repository and start exploring how the web really works, one TCP byte at a time.