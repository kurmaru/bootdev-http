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

func (handlerErr *HandlerError) Write(w io.Writer) error {
	str := fmt.Sprintf(`%v`, handlerErr.Message)
	return WriteResponse(w, []byte(str), handlerErr.Code)
}

func WriteResponse(w io.Writer, data []byte, code response.StatusCode) error {
	if err := response.WriteStatusLine(w, code); err != nil {
		return err
	}

	if err := response.WriteHeaders(w, response.GetDefaultHeaders(len(data))); err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}
