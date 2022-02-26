package container_test

import (
	"testing"

	"github.com/gabstv/container"
	"github.com/stretchr/testify/assert"
)

func TestList2D(t *testing.T) {
	// 1 0 0 0 0 0
	// 0 2 0 0 0 0
	// 0 0 3 0 0 0
	// 0 0 0 4 0 0
	// 0 0 0 0 5 0
	// 0 0 0 0 0 6
	l := container.NewList2D[int](6, 6)
	l.Set(0, 0, 1)
	l.Set(1, 1, 2)
	l.Set(2, 2, 3)
	l.Set(3, 3, 4)
	l.Set(4, 4, 5)
	l.Set(5, 5, 6)
	assert.Equal(t, 1, l.Get(0, 0))
	assert.Equal(t, 0, l.Get(1, 0))
	assert.Equal(t, 0, l.Get(2, 0))
	assert.Equal(t, 2, l.Get(1, 1))

	// 0 1 0 0 0 0
	// 0 0 2 0 0 0
	// 0 0 0 3 0 0
	// 0 0 0 0 4 0
	// 0 0 0 0 0 5
	// 6 0 0 0 0 0
	l.Shift(1, 0)
	assert.Equal(t, 0, l.Get(0, 0))
	assert.Equal(t, 1, l.Get(1, 0))
	assert.Equal(t, 0, l.Get(2, 0))
	assert.Equal(t, 0, l.Get(1, 1))
	assert.Equal(t, 2, l.Get(2, 1))

	// 0 0 0 0 0 5
	// 6 0 0 0 0 0
	// 0 1 0 0 0 0
	// 0 0 2 0 0 0
	// 0 0 0 3 0 0
	// 0 0 0 0 4 0
	l.Shift(0, 2)
	assert.Equal(t, 5, l.Get(5, 0))
	assert.Equal(t, 6, l.Get(0, 1))
	assert.Equal(t, 1, l.Get(1, 2))
	assert.Equal(t, 0, l.Get(1, 1))
	assert.Equal(t, 3, l.Get(3, 4))

	// 0 0 0 0 0 5 0
	// 6 0 0 0 0 0 0
	// 0 1 0 0 0 0 0
	// 0 0 2 0 0 0 0
	// 0 0 0 3 0 0 0
	// 0 0 0 0 4 0 0
	l.Resize(l.Width()+1, l.Height())
	l.Set(6, 1, 9)
	assert.Equal(t, 9, l.Get(6, 1))
	assert.Equal(t, 6, l.Get(0, 1))

	// 0 0 0 0
	// 6 0 0 0
	// 0 1 0 0
	l.Resize(4, 4)
	assert.Equal(t, 0, l.Get(0, 0))
	assert.Equal(t, 6, l.Get(0, 1))
	assert.Equal(t, 6, l.GetW(4, 1)) // wrapped back to 0

	l.Set(2, 1, 7)

	// 0 7
	// 1 0
	l2 := l.CopyRect(1, 1, 2, 2)
	assert.Equal(t, 0, l2.Get(0, 0))
	assert.Equal(t, 7, l2.Get(1, 0))
	assert.Equal(t, 1, l2.Get(0, 1))
	assert.Equal(t, 0, l2.Get(1, 1))
}

func TestNewList2DFrom2DSlice(t *testing.T) {
	rawl := [][]int{
		{1, 2, 3, 4},
		{4, 3, 2, 1},
		{6, 7, 8, 9, 10}, // testing out of bounds protection (this will be trimmed)
	}
	l := container.NewList2DFrom2DSlice(rawl)
	assert.Equal(t, 4, l.Width())
	assert.Equal(t, 3, l.Height())
	assert.Equal(t, 1, l.Get(0, 0))
	assert.Equal(t, 2, l.Get(1, 0))
	assert.Equal(t, 3, l.Get(1, 1))
	assert.Equal(t, 9, l.Get(3, 2))
}
