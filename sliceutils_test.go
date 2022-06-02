package container_test

import (
	"testing"

	"github.com/gabstv/container"
	"github.com/stretchr/testify/assert"
)

func TestSliceInsert(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	s = container.SliceInsert(s, 1, 10)
	assert.Equal(t, 10, s[1])
	s = container.SliceInsert(s, 5, 20)
	assert.Equal(t, 20, s[5])
}

type PlayerData struct {
	Name  string
	Score int
}

func TestSliceBinarySearch(t *testing.T) {
	s := []PlayerData{
		{Name: "a", Score: 10},
		{Name: "b", Score: 20},
		{Name: "c", Score: 30},
		{Name: "d", Score: 40},
	}

	getscore := func(v PlayerData) int {
		return v.Score
	}

	pos, exists := container.SliceBinarySearch(s, getscore, 20)
	assert.Equal(t, 1, pos)
	assert.True(t, exists)

	pos, exists = container.SliceBinarySearch(s, getscore, 50)
	assert.Equal(t, 4, pos)
	assert.False(t, exists)
}
