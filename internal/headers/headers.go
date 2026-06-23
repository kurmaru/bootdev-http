// Package headers parse/encode the binary <-> string
package headers

import (
	"bytes"
	"errors"
	"strings"
)

const (
	crlf = "\r\n"
	ws   = " "
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		return 0, true, nil
	}

	header := data[:idx]
	key, val, err := parseEntryFromString(header)
	if err != nil {
		return 0, false, err
	}

	h[key] = val

	return len(header) + len(crlf), false, nil
}

func parseEntryFromString(str []byte) (string, string, error) {
	colonIdx := bytes.Index(str, []byte(":"))

	// If colonIdx == 0 => no key => invalid
	if colonIdx <= 0 {
		return "", "", errors.New("invalid header string")
	}

	key := str[:colonIdx]
	if !validateHeaderKey(string(key)) {
		return "", "", errors.New("invalid header key")
	}

	val := strings.TrimSpace(string(str[colonIdx+1:]))

	return string(key), val, nil
}

func validateHeaderKey(str string) bool {
	for _, ch := range str {
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		default:
			if (ch >= 'A' && ch <= 'Z') ||
				(ch >= 'a' && ch <= 'z') ||
				(ch >= '0' && ch <= '9') {
				continue
			}
			return false
		}
	}
	return true
}
