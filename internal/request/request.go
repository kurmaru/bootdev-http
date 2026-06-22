// Package request handle TCP request from network
package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	parts := strings.Split(string(data), "\r\n")
	if len(parts) < 1 {
		return nil, errors.New("bad request")
	}
	reqLine, err := parseRequestLine(parts[0])
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: *reqLine,
	}, nil
}

func parseRequestLine(str string) (*RequestLine, error) {
	parts := strings.Split(str, " ")
	if len(parts) != 3 {
		return nil, errors.New("invalid request line format")
	}

	method := parts[0]
	target := parts[1]
	httpVer := parts[2]

	if method != strings.ToUpper(method) {
		return nil, errors.New("invalid method")
	}

	if httpVer != "HTTP/1.1" {
		return nil, errors.New("invalid HTTP version")
	}

	return &RequestLine{
		RequestTarget: target,
		HttpVersion:   "1.1",
		Method:        method,
	}, nil
}
