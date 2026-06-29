package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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

func handler(w response.Writer, req request.Request) {
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
