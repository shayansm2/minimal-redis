package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testInt int

func (i testInt) Gt(j testInt) bool {
	return i > j
}

func (i testInt) Eq(j testInt) bool {
	return i == j
}

func getSampleArray() []testInt {
	return []testInt{1, 10, 11, 20, 28, 42, 46, 50, 100}
}

func TestBinarySearchLessEqual(t *testing.T) {
	arr := getSampleArray()
	tests := map[testInt]int{
		1:   0,
		2:   0,
		10:  1,
		11:  2,
		12:  2,
		25:  3,
		27:  3,
		30:  4,
		42:  5,
		43:  5,
		90:  7,
		100: 8,
		200: 8,
	}
	for test, expected := range tests {
		assert.Equal(t, expected, LessEqualIndex(arr, test))
	}
}

func TestBinarySearchGreaterEqual(t *testing.T) {
	arr := getSampleArray()
	tests := map[testInt]int{
		1:   0,
		2:   1,
		10:  1,
		11:  2,
		12:  3,
		25:  4,
		27:  4,
		30:  5,
		42:  5,
		43:  6,
		90:  8,
		100: 8,
		200: 9,
	}
	for test, expected := range tests {
		assert.Equal(t, expected, GreaterEqualIndex(arr, test))
	}
}

func TestBinarySearchGreaterThan(t *testing.T) {
	arr := getSampleArray()
	tests := map[testInt]int{
		1:   1,
		2:   1,
		10:  2,
		11:  3,
		12:  3,
		25:  4,
		27:  4,
		30:  5,
		42:  6,
		43:  6,
		90:  8,
		100: 9,
		200: 9,
	}
	for test, expected := range tests {
		assert.Equal(t, expected, GreaterThanIndex(arr, test))
	}
}
