package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLinkedList(t *testing.T) {
	l := NewLinkedList[string]()
	assert.True(t, l.isEmpty())

	l.push("1")
	assert.Equal(t, "1", l.pop())
	assert.True(t, l.isEmpty())

	l.push("1")
	l.push("2")
	assert.Equal(t, "1", l.pop())
	assert.False(t, l.isEmpty())
	assert.Equal(t, "2", l.pop())
	assert.True(t, l.isEmpty())

	id := l.push("1")
	assert.False(t, l.isEmpty())
	l.del(id)
	assert.True(t, l.isEmpty())

	l.push("1")
	id = l.push("2")
	l.push("3")
	l.del(id)
	assert.Equal(t, "1", l.pop())
	assert.False(t, l.isEmpty())
	assert.Equal(t, "3", l.pop())
	assert.True(t, l.isEmpty())

	id = l.push("1")
	l.push("2")
	l.push("3")
	l.del(id)
	assert.Equal(t, "2", l.pop())
	assert.False(t, l.isEmpty())
	assert.Equal(t, "3", l.pop())
	assert.True(t, l.isEmpty())

	l.push("1")
	l.push("2")
	id = l.push("3")
	l.del(id)
	assert.Equal(t, "1", l.pop())
	assert.False(t, l.isEmpty())
	assert.Equal(t, "2", l.pop())
	assert.True(t, l.isEmpty())

	l.push("1")
	l.push("2")
	id = l.push("3")
	l.del(id)
	l.push("4")
	assert.Equal(t, "1", l.pop())
	assert.Equal(t, "2", l.pop())
	assert.Equal(t, "4", l.pop())
}
