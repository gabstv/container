package container_test

import (
	"testing"

	"github.com/gabstv/container"
	"github.com/stretchr/testify/assert"
)

func TestSortMap(t *testing.T) {
	m := make(map[string]int)
	m["a"] = 1
	m["b"] = 2
	m["c"] = 3
	m["d"] = 4
	m["e"] = -1

	slc := container.MapToSlice(m)
	container.Sort(slc, func(a, b container.MI[string, int]) bool {
		return a.Value > b.Value
	})
	assert.Equal(t, "d", slc[0].Key)
	assert.Equal(t, 3, slc[1].Value)

	m2 := container.SliceToMap(slc[:2])
	assert.Equal(t, 2, len(m2))
	assert.Equal(t, 0, m2["a"]) // not found; returns zero value
	assert.Equal(t, 4, m2["d"])
}
