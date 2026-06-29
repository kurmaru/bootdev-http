// Package response handle sending the HTTP response back to client
package response

import (
	"strconv"

	"github.com/kurmaru/bootdev-http/internal/headers"
)

type StatusCode int

const (
	OK                  StatusCode = 200
	BadRequest          StatusCode = 400
	InternalServerError StatusCode = 500
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("connection", "close")
	h.Set("content-type", "text/plain")
	h.Set("Content-Length", strconv.Itoa(contentLen))
	return h
}
