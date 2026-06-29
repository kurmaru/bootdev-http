package server

import (
	"fmt"
	"io"

	"github.com/kurmaru/bootdev-http/internal/request"
	"github.com/kurmaru/bootdev-http/internal/response"
)

type HandlerError struct {
	Code    response.StatusCode
	Message string
}

type Handler func(w io.Writer, req request.Request) *HandlerError

func WriteError(w io.Writer, handlerErr *HandlerError) error {
	if err := response.WriteStatusLine(w, handlerErr.Code); err != nil {
		return err
	}
	str := fmt.Sprintf(`%v`, handlerErr.Message)
	if err := response.WriteHeaders(w, response.GetDefaultHeaders(len(str))); err != nil {
		return err
	}

	err := response.WriteBody(w, []byte(str))

	return err
}
