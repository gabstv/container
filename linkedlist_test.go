package container_test

import (
	"testing"

	"github.com/gabstv/container"
	"github.com/stretchr/testify/assert"
)

func TestDoubleLinkedList(t *testing.T) {
	var ll container.LinkedList[int]
	ll.Push(100)
	ll.Push(200)
	item3 := ll.Push(300)
	item4 := ll.Push(400)
	assert.True(t, item4.Remove())
	assert.False(t, item4.Remove())
	assert.Equal(t, item3, ll.Last())
	assert.Equal(t, item3.Data(), ll.Last().Data())
	assert.Equal(t, 100, ll.First().Data())
	assert.Equal(t, 200, ll.First().Next().Data())
	assert.Equal(t, 300, ll.First().Next().Next().Data())
	assert.Equal(t, 300, ll.Pop())
	assert.Equal(t, 2, ll.Len())
	assert.Equal(t, 200, ll.Pop())
	assert.Equal(t, 100, ll.Pop())
	assert.Equal(t, 0, ll.Pop()) // zero value
	ll.Push(150)
	ll.Push(250)
	ll.Unshift(50)
	assert.Equal(t, 3, ll.Len())
	assert.Equal(t, 50, ll.Shift())
	assert.Equal(t, 150, ll.Shift())
	assert.Equal(t, 250, ll.Shift())
	assert.Equal(t, 0, ll.Len())
	ll.Push(123)
	ll.Push(456)
	ll.Push(789)
	ll.Remove(456)
	assert.Equal(t, 2, ll.Len())
	assert.Equal(t, 123, ll.First().Data())
	assert.Equal(t, 789, ll.Last().Data())
	assert.True(t, ll.Contains(789))
	assert.False(t, ll.Contains(999))
	assert.Equal(t, 123, ll.Shift())
	assert.Equal(t, 789, ll.Shift())
	assert.Equal(t, 0, ll.Shift()) // zero value
	assert.False(t, ll.Remove(-1))
	ll.Push(123)
	ll.Push(124)
	ll.Push(125)
	assert.Equal(t, 124, ll.Last().Prev().Data())
	ll.Clear()
	assert.Equal(t, 0, ll.Pop()) // zero value
	ll.Clear()
	assert.Equal(t, 1234, ll.Unshift(1234).Data())
	assert.Nil(t, ll.First().Prev().Prev())
}
