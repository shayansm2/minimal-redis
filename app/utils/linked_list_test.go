package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyLinkedList(t *testing.T) {
	l := NewLinkedList[string]()
	assert.True(t, l.IsEmpty())
}

func TestLinkedListPushPop(t *testing.T) {
	l := NewLinkedList[string]()

	l.Push("1")
	assert.Equal(t, "1", l.Pop())
	assert.True(t, l.IsEmpty())

	l.Push("1")
	l.Push("2")
	assert.Equal(t, "1", l.Pop())
	assert.False(t, l.IsEmpty())
	assert.Equal(t, "2", l.Pop())
	assert.True(t, l.IsEmpty())
}

func TestLinkedListDelete(t *testing.T) {
	l := NewLinkedList[string]()

	id := l.Push("1")
	assert.False(t, l.IsEmpty())
	l.Del(id)
	assert.True(t, l.IsEmpty())

	l.Push("1")
	id = l.Push("2")
	l.Push("3")
	l.Del(id)
	assert.Equal(t, "1", l.Pop())
	assert.False(t, l.IsEmpty())
	assert.Equal(t, "3", l.Pop())
	assert.True(t, l.IsEmpty())

	id = l.Push("1")
	l.Push("2")
	l.Push("3")
	l.Del(id)
	assert.Equal(t, "2", l.Pop())
	assert.False(t, l.IsEmpty())
	assert.Equal(t, "3", l.Pop())
	assert.True(t, l.IsEmpty())

	l.Push("1")
	l.Push("2")
	id = l.Push("3")
	l.Del(id)
	assert.Equal(t, "1", l.Pop())
	assert.False(t, l.IsEmpty())
	assert.Equal(t, "2", l.Pop())
	assert.True(t, l.IsEmpty())

	l.Push("1")
	l.Push("2")
	id = l.Push("3")
	l.Del(id)
	l.Push("4")
	assert.Equal(t, "1", l.Pop())
	assert.Equal(t, "2", l.Pop())
	assert.Equal(t, "4", l.Pop())
}
