package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBulkEncode(t *testing.T) {
	raw := "hello"
	expected := "$5\r\nhello\r\n"
	assert.Equal(t, bulkEncode(raw), expected)

	raw = ""
	expected = "$0\r\n\r\n"
	assert.Equal(t, bulkEncode(raw), expected)
}

func TestBulkDecode(t *testing.T) {
	expected := "hello"
	encoded := "$5\r\nhello\r\n"
	decoded, err := bulkDecode(encoded)
	assert.Nil(t, err)
	assert.Equal(t, expected, decoded)

	expected = ""
	encoded = "$0\r\n\r\n"
	decoded, err = bulkDecode(encoded)
	assert.Nil(t, err)
	assert.Equal(t, expected, decoded)
}
