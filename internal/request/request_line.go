package request

import (
	"bytes"
	"errors"
	"strings"
)

func parseRequestLine(str []byte) (rq *RequestLine, consumed int, err error) {
	idx := bytes.Index([]byte(str), []byte(crlf))
	if idx == -1 {
		return nil, 0, nil
	}

	return requestLineFromString(string(str[:idx]))
}

func requestLineFromString(str string) (rq *RequestLine, consumed int, err error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, 0, errors.New("invalid request line format")
	}

	method := parts[0]
	target := parts[1]
	httpVer := parts[2]

	if method != strings.ToUpper(method) {
		return nil, 0, errors.New("invalid method")
	}

	if httpVer != "HTTP/1.1" {
		return nil, 0, errors.New("invalid HTTP version")
	}

	return &RequestLine{
		RequestTarget: target,
		HttpVersion:   "1.1",
		Method:        method,
	}, len(str) + 2, nil
}
