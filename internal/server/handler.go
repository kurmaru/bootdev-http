package server

import (
	"github.com/kurmaru/bootdev-http/internal/request"
	"github.com/kurmaru/bootdev-http/internal/response"
)

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

type Handler func(w response.Writer, req request.Request)
