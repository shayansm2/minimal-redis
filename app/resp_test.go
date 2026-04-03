package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespArrayParse(t *testing.T) {
	encoded := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	expected := []string{"hello", "world"}
	decoded, err := respArrayParse(encoded)
	assert.Nil(t, err)
	assert.Equal(t, expected, decoded)

	encoded = "*0\r\n"
	expected = []string{}
	decoded, err = respArrayParse(encoded)
	assert.Nil(t, err)
	assert.Equal(t, expected, decoded)
}

func TestBulkStringEncode(t *testing.T) {
	raw := "hello"
	expected := "$5\r\nhello\r\n"
	assert.Equal(t, toBulkString(raw), expected)

	raw = ""
	expected = "$0\r\n\r\n"
	assert.Equal(t, toBulkString(raw), expected)
}

func TestBulkStringDecode(t *testing.T) {
	expected := "hello"
	encoded := "$5\r\nhello\r\n"
	decoded, err := bulkStringDecode(encoded)
	assert.Nil(t, err)
	assert.Equal(t, expected, decoded)

	expected = ""
	encoded = "$0\r\n\r\n"
	decoded, err = bulkStringDecode(encoded)
	assert.Nil(t, err)
	assert.Equal(t, expected, decoded)
}

func TestRespArrayEncode(t *testing.T) {
	expected := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	array := []string{"hello", "world"}
	encoded := toRespArray(array)
	assert.Equal(t, expected, encoded)

	expected = "*0\r\n"
	array = []string{}
	encoded = toRespArray(array)
	assert.Equal(t, expected, encoded)
}
