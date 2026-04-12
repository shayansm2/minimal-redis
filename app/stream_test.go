package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getSampleStream() *Stream {
	return &Stream{
		Entry{id: ID{0, 1}},  // 0
		Entry{id: ID{1, 0}},  // 1
		Entry{id: ID{1, 1}},  // 2
		Entry{id: ID{2, 0}},  // 3
		Entry{id: ID{2, 8}},  // 4
		Entry{id: ID{4, 2}},  // 5
		Entry{id: ID{4, 6}},  // 6
		Entry{id: ID{5, 0}},  // 7
		Entry{id: ID{10, 0}}, // 8
	}
}
func TestBinarySearchLessEqual(t *testing.T) {
	stream := getSampleStream()
	assert.Equal(t, 0, leIdx(stream, ID{0, 1}))
	assert.Equal(t, 0, leIdx(stream, ID{0, 2}))
	assert.Equal(t, 1, leIdx(stream, ID{1, 0}))
	assert.Equal(t, 2, leIdx(stream, ID{1, 1}))
	assert.Equal(t, 2, leIdx(stream, ID{1, 2}))
	assert.Equal(t, 3, leIdx(stream, ID{2, 5}))
	assert.Equal(t, 3, leIdx(stream, ID{2, 7}))
	assert.Equal(t, 4, leIdx(stream, ID{3, 0}))
	assert.Equal(t, 5, leIdx(stream, ID{4, 2}))
	assert.Equal(t, 5, leIdx(stream, ID{4, 3}))
	assert.Equal(t, 7, leIdx(stream, ID{9, 0}))
	assert.Equal(t, 8, leIdx(stream, ID{10, 0}))
	assert.Equal(t, 8, leIdx(stream, ID{20, 0}))
}

func TestBinarySearchGreaterEqual(t *testing.T) {
	stream := getSampleStream()
	assert.Equal(t, 0, geIdx(stream, ID{0, 1}))
	assert.Equal(t, 1, geIdx(stream, ID{0, 2}))
	assert.Equal(t, 1, geIdx(stream, ID{1, 0}))
	assert.Equal(t, 2, geIdx(stream, ID{1, 1}))
	assert.Equal(t, 3, geIdx(stream, ID{1, 2}))
	assert.Equal(t, 4, geIdx(stream, ID{2, 5}))
	assert.Equal(t, 4, geIdx(stream, ID{2, 7}))
	assert.Equal(t, 5, geIdx(stream, ID{3, 0}))
	assert.Equal(t, 5, geIdx(stream, ID{4, 2}))
	assert.Equal(t, 6, geIdx(stream, ID{4, 3}))
	assert.Equal(t, 8, geIdx(stream, ID{9, 0}))
	assert.Equal(t, 8, geIdx(stream, ID{10, 0}))
	assert.Equal(t, 9, geIdx(stream, ID{20, 0}))
}
