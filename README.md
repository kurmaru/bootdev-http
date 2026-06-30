# HTTP from TCP

Building HTTP/1.1 from scratch on top of raw TCP sockets — no `net/http` used.

Created as part of the [Boot.dev](https://boot.dev) "HTTP from TCP" course.

## Project Structure

```
├── cmd/
│   ├── httpserver/       Production HTTP server (port 42069)
│   ├── tcplistener/      Debug tool — prints parsed HTTP requests
│   └── udpsender/        UDP client — sends stdin to UDP port 42069
└── internal/
    ├── headers/          HTTP header parsing (RFC-compliant)
    ├── request/          HTTP request parser (state machine)
    ├── response/         HTTP response writer + chunked encoding
    └── server/           Concurrent TCP server framework
```

## Features

- Full HTTP/1.1 request parsing via state machine (variable-length reads)
- RFC-compliant header parsing (case-insensitive keys, duplicate combining)
- Chunked transfer encoding with trailers (SHA256, content-length)
- Concurrent connection handling (goroutine-per-connection)
- Routes: static files, error codes, reverse proxy to httpbin.org
- Zero external dependencies beyond testify (tests only)

## Run

```bash
# Start the server
go run ./cmd/httpserver/main.go

# Test it
curl http://localhost:42069/
curl http://localhost:42069/video
```

## Test

```bash
go test -v ./...
```
