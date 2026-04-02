package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespParse(t *testing.T) {
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
