package response

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kurmaru/bootdev-http/internal/headers"
)

var crlf = []byte("\r\n")

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	size := len(p)
	hex := strconv.FormatInt(int64(size), 16)
	count := 0

	writeCount, err := w.writer.Write([]byte(hex))
	if err != nil {
		return count, err
	}
	count += writeCount

	writeCount, err = w.writer.Write(crlf)
	if err != nil {
		return count, err
	}
	count += writeCount

	writeCount, err = w.writer.Write(p)
	if err != nil {
		return count, err
	}
	count += writeCount

	writeCount, err = w.writer.Write(crlf)
	if err != nil {
		return count, err
	}
	count += writeCount

	return count, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.writer.Write([]byte("0\r\n\r\n"))
}

func (w *Writer) WriteTrailer(h headers.Headers) error {
	if _, err := w.writer.Write([]byte("0\r\n")); err != nil {
		return err
	}

	var str strings.Builder
	for key, val := range h {
		fmt.Fprintf(&str, "%v: %v\r\n", key, val)
	}
	_, err := str.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	_, err = w.writer.Write([]byte(str.String()))
	if err != nil {
		return err
	}

	return nil
}
