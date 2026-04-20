package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRespArrayParse(t *testing.T) {
	tests := []struct {
		name    string
		encoded string
		decoded [][]string
	}{
		{
			name:    "test resp array parse",
			encoded: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			decoded: [][]string{{"hello", "world"}},
		},
		{
			name:    "test null array parse",
			encoded: "*0\r\n",
			decoded: [][]string{{}},
		},
		{
			name:    "test array parse containing $ character",
			encoded: "*6\r\n$5\r\nXREAD\r\n$5\r\nBLOCK\r\n$4\r\n1000\r\n$7\r\nstreams\r\n$8\r\nsome_key\r\n$1\r\n$",
			decoded: [][]string{{"XREAD", "BLOCK", "1000", "streams", "some_key", "$"}},
		},
		{
			name:    "test multi array parse",
			encoded: "*2\r\n$6\r\nSELECT\r\n$1\r\n0\r\n*3\r\n$3\r\nSET\r\n$3\r\nfoo\r\n$3\r\n123\r\n",
			decoded: [][]string{{"SELECT", "0"}, {"SET", "foo", "123"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := make([][]string, 0)
			for encoded := tt.encoded; encoded != ""; {
				decoded, n, err := respArrayBulkStringParse(encoded)
				assert.Nil(t, err)
				actual = append(actual, decoded)
				encoded = encoded[n:]
			}
			assert.Equal(t, tt.decoded, actual)
		})
	}
}

func TestBulkStringEncode(t *testing.T) {
	raw := "hello"
	expected := "$5\r\nhello\r\n"
	assert.Equal(t, toBulkString(BulkStr(raw)), expected)

	raw = ""
	expected = "$0\r\n\r\n"
	assert.Equal(t, toBulkString(BulkStr(raw)), expected)
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

func TestRespStringDecode(t *testing.T) {
	decoded, err := respSimpleStringDecode("+PONG\r\n")
	assert.Nil(t, err)
	assert.Equal(t, "PONG", decoded)
}

func TestRespArrayEncode(t *testing.T) {
	expected := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	bulkArray := []BulkStr{"hello", "world"}
	array := make([]string, len(bulkArray))
	for i, bulkStr := range bulkArray {
		array[i], _ = encode(bulkStr)
	}
	assert.Equal(t, expected, toRespArray(array))

	expected = "*0\r\n"
	array = []string{}
	assert.Equal(t, expected, toRespArray(array))

	expected = "*3\r\n:1\r\n:2\r\n:3\r\n"
	intArray := []int{1, 2, 3}
	array = make([]string, len(intArray))
	for i, n := range intArray {
		array[i], _ = encode(n)
	}
	assert.Equal(t, expected, toRespArray(array))
}

func TestEncode(t *testing.T) {
	val, err := encode(nil)
	assert.Nil(t, err)
	assert.Equal(t, NullBulkString, val)

	val, err = encode(22)
	assert.Nil(t, err)
	assert.Equal(t, ":22\r\n", val)

	val, err = encode(RespStr("OK"))
	assert.Nil(t, err)
	assert.Equal(t, "+OK\r\n", val)

	val, err = encode(BulkStr("hi"))
	assert.Nil(t, err)
	assert.Equal(t, "$2\r\nhi\r\n", val)
}

func TestArrayEncode(t *testing.T) {
	var nullArr []string
	val, err := encode(nullArr)
	assert.Equal(t, NullArray, val)

	val, err = encode([]BulkStr{"PING"})
	assert.Nil(t, err)
	assert.Equal(t, "*1\r\n$4\r\nPING\r\n", val)

	expected := "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	encoded, err := encode([]BulkStr{"hello", "world"})
	assert.Nil(t, err)
	assert.Equal(t, expected, encoded)

	encoded, err = encode([]string{})
	assert.Nil(t, err)
	assert.Equal(t, "*0\r\n", encoded)

	expected = "*3\r\n:1\r\n:2\r\n:3\r\n"
	encoded, err = encode([]int{1, 2, 3})
	assert.Nil(t, err)
	assert.Equal(t, expected, encoded)

	raw := []interface{}{
		[]interface{}{
			BulkStr("0-2"),
			[]BulkStr{"bar", "baz"},
		},
		[]interface{}{
			BulkStr("0-3"),
			[]BulkStr{"baz", "foo"},
		},
	}
	expected = "*2\r\n*2\r\n$3\r\n0-2\r\n*2\r\n$3\r\nbar\r\n$3\r\nbaz\r\n*2\r\n$3\r\n0-3\r\n*2\r\n$3\r\nbaz\r\n$3\r\nfoo\r\n"
	val, err = encode(raw)
	assert.Nil(t, err)
	assert.Equal(t, expected, val)
}
