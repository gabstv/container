package container_test

import (
	"testing"

	"github.com/gabstv/container"
	"github.com/stretchr/testify/assert"
)

func TestSortedDictionary(t *testing.T) {
	var d container.SortedDictionary[int64, string]
	d.Set(100, "one")
	d.Set(200, "two")
	d.Set(150, "three")
	assert.Equal(t, "one", d.Get(100))
	assert.Equal(t, "two", d.Get(200))
	assert.Equal(t, "three", d.Get(150))
	assert.Equal(t, 1, d.Index(150))
	assert.Equal(t, 2, d.Index(200))
}
