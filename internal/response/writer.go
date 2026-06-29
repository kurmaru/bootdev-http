package response

import (
	"fmt"
	"io"
	"strings"

	"github.com/kurmaru/bootdev-http/internal/headers"
)

type WriterState string

const (
	statusLineState WriterState = "status line"
	headersState    WriterState = "headers"
	bodyState       WriterState = "body"
	trailersState   WriterState = "trailers"
	doneState       WriterState = "done"
)

type Writer struct {
	writer io.Writer
	state  WriterState
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{writer: w, state: statusLineState}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != statusLineState {
		return fmt.Errorf("invalid state - expect: %v - got %v", statusLineState, w.state)
	}

	reason := ""
	switch statusCode {
	case OK:
		reason = "OK"
	case BadRequest:
		reason = "Bad Request"
	case InternalServerError:
		reason = "Internal Server Error"
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %v %v\r\n", statusCode, reason)
	if _, err := w.writer.Write([]byte(statusLine)); err != nil {
		return err
	}

	w.state = headersState
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.state != headersState {
		return fmt.Errorf("invalid state - expect: %v - got %v", headersState, w.state)
	}

	var str strings.Builder

	for key, val := range headers {
		fmt.Fprintf(&str, "%v: %v\r\n", key, val)
	}
	str.Write([]byte("\r\n"))

	_, err := w.writer.Write([]byte(str.String()))
	if err != nil {
		return err
	}

	w.state = bodyState
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != bodyState {
		return 0, fmt.Errorf("invalid state - expect: %v - got %v", bodyState, w.state)
	}

	return w.writer.Write(p)
}
