package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid single header
	headers = NewHeaders()
	data = []byte("hOsT: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["hOsT"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("Host : localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}

func TestHeadersKey(t *testing.T) {
	// Test: Valid special char
	headers := NewHeaders()
	data := []byte("H&st: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["H&st"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid special char
	headers = NewHeaders()
	data = []byte("H$s!!~~~t: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["H$s!!~~~t"])
	assert.Equal(t, 28, n)
	assert.False(t, done)

	// Test: Invalid special char
	headers = NewHeaders()
	data = []byte("H@s!!\\~~~t: localhost:42069\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.False(t, done)

	// Test: Invalid unicode character
	headers = NewHeaders()
	data = []byte("Hót: localhost:42069\r\n\r\n")
	_, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.False(t, done)
}
