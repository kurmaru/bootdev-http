// Package request handle TCP request from network
package request

import (
	"bytes"
	"errors"
	"io"
	"strings"
)

const (
	crlf        = "\r\n"
	readBufSize = 8
)

type RequestState int

const (
	Initialized RequestState = iota
	Done
)

type Request struct {
	RequestLine RequestLine
	State       RequestState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	curBuf := make([]byte, readBufSize)
	readToIdx := 0
	req := Request{
		State: Initialized,
	}

	for {
		if readToIdx == cap(curBuf) {
			buf := make([]byte, cap(curBuf)*2)
			copy(buf, curBuf)
			curBuf = buf
		}

		count, err := reader.Read(curBuf[readToIdx:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				req.State = Done
			} else {
				return nil, err
			}
		}
		readToIdx += count

		parsedCount, err := req.parse(curBuf)
		if err != nil {
			return nil, err
		}

		if parsedCount > 0 {
			readToIdx -= parsedCount
			buf := make([]byte, readToIdx)
			copy(buf, curBuf[parsedCount:])
			curBuf = buf
		}

		if req.State == Done {
			break
		}
	}

	return &req, nil
}

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
	}, len(str), nil
}

func (r *Request) parse(data []byte) (int, error) {
	switch r.State {
	case Initialized:
		rq, count, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if count == 0 {
			return 0, nil
		}
		r.RequestLine = *rq
		r.State = Done
		return count, nil
	case Done:
		return 0, errors.New("trying to read data from done state")
	default:
		return 0, errors.New("unknown state")
	}
}
