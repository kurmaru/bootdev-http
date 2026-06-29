package response

import (
	"strconv"
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
