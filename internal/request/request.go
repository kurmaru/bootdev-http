// Package request handle TCP request from network
package request

import (
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/kurmaru/bootdev-http/internal/headers"
)

const (
	crlf        = "\r\n"
	readBufSize = 8
)

type RequestState int

const (
	requestStateInitialized RequestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

type Request struct {
	State          RequestState
	RequestLine    RequestLine
	Headers        headers.Headers
	Body           []byte
	bodyLengthRead int
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
		State:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}

	for req.State != requestStateDone {
		if readToIdx == cap(curBuf) {
			buf := make([]byte, cap(curBuf)*2)
			copy(buf, curBuf)
			curBuf = buf
		}

		count, err := reader.Read(curBuf[readToIdx:])
		if err != nil {
			if errors.Is(err, io.EOF) {
				if req.State != requestStateDone {
					return nil, errors.New("get EOF in middle of parsing")
				}
				break
			} else {
				return nil, err
			}
		}
		readToIdx += count

		parsedCount, err := req.parse(curBuf[:readToIdx])
		if err != nil {
			return nil, err
		}

		// If parsed success, trim down the size of buf
		if parsedCount > 0 {
			buf := make([]byte, max(len(curBuf), readBufSize))
			copy(buf, curBuf[parsedCount:])
			curBuf = buf
			readToIdx -= parsedCount
		}

	}

	return &req, nil
}

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0
	for r.State != requestStateDone {
		parsedCount, err := r.parseSingleLine(data[totalBytesParsed:])
		if err != nil || parsedCount == 0 {
			return totalBytesParsed, err
		}
		totalBytesParsed += parsedCount
	}
	return totalBytesParsed, nil
}

func (r *Request) parseSingleLine(data []byte) (int, error) {
	switch r.State {
	case requestStateInitialized:
		rq, count, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if count == 0 {
			return 0, nil
		}
		r.RequestLine = *rq
		r.State = requestStateParsingHeaders
		return count, nil
	case requestStateParsingHeaders:
		count, done, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}
		if done {
			r.State = requestStateParsingBody
		}
		return count, nil
	case requestStateParsingBody:
		contentLenStr, ok := r.Headers.Get("Content-Length")
		if !ok {
			// assume that if no content-length header is present, there is no body
			r.State = requestStateDone
			return len(data), nil
		}
		contentLen, err := strconv.Atoi(contentLenStr)
		if err != nil {
			return 0, fmt.Errorf("malformed Content-Length: %s", err)
		}
		r.Body = append(r.Body, data...)
		r.bodyLengthRead += len(data)
		if r.bodyLengthRead > contentLen {
			return 0, fmt.Errorf("Content-Length too large")
		}
		if r.bodyLengthRead == contentLen {
			r.State = requestStateDone
		}
		return len(data), nil
	case requestStateDone:
		return 0, errors.New("trying to read data from done state")
	default:
		return 0, errors.New("unknown state")
	}
}
