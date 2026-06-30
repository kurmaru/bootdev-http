package main

import (
	"crypto/sha256"
	"encoding/hex"
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

func internalHandlers(w response.Writer, r request.Request) {
	if r.RequestLine.RequestTarget == "/video" && r.RequestLine.Method == "GET" {
		videoHandler(w, r)
		return
	}
	htmlHandler(w, r)
}

func htmlHandler(w response.Writer, r request.Request) {
	statusCode := response.OK
	var content []byte
	h := response.GetDefaultHeaders(0)

	switch r.RequestLine.RequestTarget {
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

	h.WriteHeaders("Content-Length", strconv.Itoa(len(content)))
	h.WriteHeaders("Content-Type", "text/html")

	w.WriteStatusLine(statusCode)
	w.WriteHeaders(h)
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

	h := headers.NewHeaders()
	for key, vals := range proxReq.Header {
		if len(vals) > 0 {
			h.WriteHeaders(key, vals[0])
		}
	}
	h.Delete("Content-Length")
	h.Set("connection", "close")
	h.WriteHeaders("Transfer-Encoding", "chunked")
	h.WriteHeaders("Trailer", "X-Content-SHA256")
	h.Set("Trailer", "X-Content-Length")

	body := make([]byte, 0, 1024)

	if err := w.WriteHeaders(h); err != nil {
		fmt.Println(err)
		return
	}

	buf := make([]byte, 1024)

	for {
		n, err := proxReq.Body.Read(buf)
		if n > 0 {
			fmt.Printf("Read %v bytes\n", n)
			body = append(body, buf[:n]...)
			w.WriteChunkedBody(buf[:n])
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Printf("Failed to read response: %v", err)
			return
		}
	}

	// Body have a whitespace in the end
	body = body[:len(body)-1]
	hash := sha256.Sum256(body)

	trailers := headers.NewHeaders()
	trailers.WriteHeaders("X-Content-SHA256", hex.EncodeToString(hash[:]))
	trailers.WriteHeaders("X-Content-Length", strconv.Itoa(len(body)))

	if err := w.WriteTrailer(trailers); err != nil {
		fmt.Printf("Failed to write trailers: %v\n", err)
	}
}

func videoHandler(w response.Writer, r request.Request) {
	file, err := os.Open("assets/vim.mp4")
	if err != nil {
		body := fmt.Sprintf("Failed to open asset: %v", err)
		w.WriteStatusLine(response.InternalServerError)
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody([]byte(body))
		return
	}
	defer file.Close()
	buf, err := io.ReadAll(file)
	if err != nil {
		fmt.Printf("Failed to read file: %v", err)
		return
	}

	w.WriteStatusLine(response.OK)

	h := response.GetDefaultHeaders(len(buf))
	h.WriteHeaders("Content-Type", "video/mp4")
	w.WriteHeaders(h)

	w.WriteBody(buf)
}
