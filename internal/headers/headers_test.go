package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParser(t *testing.T) {
	// Test: Valid single header
	headers := make(Headers)
	data := []byte("host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = make(Headers)
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid character in header
	headers = make(Headers)
	data = []byte("       H(st : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Uppercase FieldName should add as a lowercase key
	headers = make(Headers)
	data = []byte("Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: multipe same fieldname should have values comma separated
	headers = make(Headers)
	data = []byte("Host: localhost:42069\r\nHost: localhost:42070\r\n\r\n")
	bytesConsumed0, _, _ := headers.Parse(data)
	bytesConsumed1, done, err := headers.Parse(data[bytesConsumed0:])
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069,localhost:42070", headers["host"])
	assert.Equal(t, bytesConsumed0+bytesConsumed1, len(data)-2)
	assert.False(t, done)
}
