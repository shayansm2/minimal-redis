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
	tests := []struct {
		id  ID
		idx int
	}{
		{ID{0, 1}, 0},
		{ID{0, 2}, 0},
		{ID{1, 0}, 1},
		{ID{1, 1}, 2},
		{ID{1, 2}, 2},
		{ID{2, 5}, 3},
		{ID{2, 7}, 3},
		{ID{3, 0}, 4},
		{ID{4, 2}, 5},
		{ID{4, 3}, 5},
		{ID{9, 0}, 7},
		{ID{10, 0}, 8},
		{ID{20, 0}, 8},
	}
	for _, test := range tests {
		assert.Equal(t, test.idx, leIdx(stream, test.id))
	}
}

func TestBinarySearchGreaterEqual(t *testing.T) {
	stream := getSampleStream()
	tests := []struct {
		id  ID
		idx int
	}{
		{ID{0, 1}, 0},
		{ID{0, 2}, 1},
		{ID{1, 0}, 1},
		{ID{1, 1}, 2},
		{ID{1, 2}, 3},
		{ID{2, 5}, 4},
		{ID{2, 7}, 4},
		{ID{3, 0}, 5},
		{ID{4, 2}, 5},
		{ID{4, 3}, 6},
		{ID{9, 0}, 8},
		{ID{10, 0}, 8},
		{ID{20, 0}, 9},
	}
	for _, test := range tests {
		assert.Equal(t, test.idx, geIdx(stream, test.id))
	}
}

func TestBinarySearchGreaterThan(t *testing.T) {
	stream := getSampleStream()
	tests := []struct {
		id  ID
		idx int
	}{
		{ID{0, 1}, 1},
		{ID{0, 2}, 1},
		{ID{1, 0}, 2},
		{ID{1, 1}, 3},
		{ID{1, 2}, 3},
		{ID{2, 5}, 4},
		{ID{2, 7}, 4},
		{ID{3, 0}, 5},
		{ID{4, 2}, 6},
		{ID{4, 3}, 6},
		{ID{9, 0}, 8},
		{ID{10, 0}, 9},
		{ID{20, 0}, 9},
	}
	for _, test := range tests {
		assert.Equal(t, test.idx, gdIdx(stream, test.id))
	}
}
