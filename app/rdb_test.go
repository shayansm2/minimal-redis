package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLengthEncoding(t *testing.T) {
	v, _, valType := getLen([]byte{0x0A})
	assert.Equal(t, 10, v)
	assert.Equal(t, ValueTypeStr, valType)

	v, _, valType = getLen([]byte{0x42, 0xBC})
	assert.Equal(t, 700, v)
	assert.Equal(t, ValueTypeStr, valType)

	v, _, valType = getLen([]byte{0x80, 0x00, 0x00, 0x42, 0x68})
	assert.Equal(t, 17000, v)
	assert.Equal(t, ValueTypeStr, valType)

	v, _, valType = getLen([]byte{0xC0})
	assert.Equal(t, 1, v)
	assert.Equal(t, ValueTypeInt, valType)

	v, _, valType = getLen([]byte{0xC1})
	assert.Equal(t, 2, v)
	assert.Equal(t, ValueTypeInt, valType)

	v, _, valType = getLen([]byte{0xC2})
	assert.Equal(t, 4, v)
	assert.Equal(t, ValueTypeInt, valType)
}
