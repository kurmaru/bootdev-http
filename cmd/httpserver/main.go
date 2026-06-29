package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/kurmaru/bootdev-http/internal/headers"
	"github.com/kurmaru/bootdev-http/internal/request"
	"github.com/kurmaru/bootdev-http/internal/response"
	"github.com/kurmaru/bootdev-http/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

const (
	httpbinPath    = "/httpbin"
	httpbinRootURL = "https://httpbin.org"
)

func handler(w response.Writer, req request.Request) {
	path := req.RequestLine.RequestTarget
	if strings.HasPrefix(path, httpbinPath) {
		req.RequestLine.RequestTarget = strings.TrimPrefix(path, httpbinPath)
		proxyHandler(w, req)
		return
	}

	internalHandlers(w, req)
}

func internalHandlers(w response.Writer, req request.Request) {
	statusCode := response.OK
	var content []byte
	headers := response.GetDefaultHeaders(0)

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		statusCode = response.BadRequest
		content = []byte(`
		<html>
			<head>
				<title>400 Bad Request</title>
			</head>
			<body>
				<h1>Bad Request</h1>
				<p>Your request honestly kinda sucked.</p>
			</body>
		</html>
		`)
	case "/myproblem":
		statusCode = response.InternalServerError
		content = []byte(`
		<html>
			<head>
				<title>500 Internal Server Error</title>
			</head>
			<body>
				<h1>Internal Server Error</h1>
				<p>Okay, you know what? This one is on me.</p>
			</body>
		</html>
		`)
	default:
		content = []byte(`
		<html>
			<head>
				<title>200 OK</title>
			</head>
			<body>
				<h1>Success!</h1>
				<p>Your request was an absolute banger.</p>
			</body>
		</html>
		`)
	}

	headers.WriteHeaders("Content-Length", strconv.Itoa(len(content)))
	headers.WriteHeaders("Content-Type", "text/html")

	w.WriteStatusLine(statusCode)
	w.WriteHeaders(headers)
	w.WriteBody(content)
}

func proxyHandler(w response.Writer, req request.Request) {
	proxReq, err := http.Get(httpbinRootURL + req.RequestLine.RequestTarget)
	if err != nil {
		body := fmt.Sprintf("Failed to connect to proxy: %v", err)
		w.WriteStatusLine(response.BadRequest)
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody([]byte(body))
		return
	}
	defer proxReq.Body.Close()

	if err := w.WriteStatusLine(response.StatusCode(proxReq.StatusCode)); err != nil {
		fmt.Println(err)
		return
	}

	headers := headers.NewHeaders()
	for h, vals := range proxReq.Header {
		if len(vals) > 0 {
			headers.WriteHeaders(h, vals[0])
		}
	}
	headers.Delete("Content-Length")
	headers.WriteHeaders("Transfer-Encoding", "chunked")

	if err := w.WriteHeaders(headers); err != nil {
		fmt.Println(err)
		return
	}

	buf := make([]byte, 1024)

	for {
		n, err := proxReq.Body.Read(buf)
		if n > 0 {
			fmt.Printf("Read %v bytes\n", n)
			w.WriteChunkedBody(buf[:n])
		}

		if err != nil {
			if err == io.EOF {
				break
			}

			body := fmt.Sprintf("Failed to read response: %v", err)
			w.WriteStatusLine(response.InternalServerError)
			w.WriteHeaders(response.GetDefaultHeaders(len(body)))
			w.WriteBody([]byte(body))
			return
		}
	}

	if _, err := w.WriteChunkedBodyDone(); err != nil {
		fmt.Printf("Failed to write chunk done: %v\n", err)
	}
}
