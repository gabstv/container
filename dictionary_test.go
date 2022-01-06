package container_test

import (
	"testing"

	"github.com/gabstv/container"
	"github.com/stretchr/testify/assert"
)

type TestKey struct {
	Name  string
	Score int
}

func TestDictionary(t *testing.T) {
	var d container.Dictionary[int64, string]
	d.Set(1, "one")
	d.Set(2, "two")
	assert.Equal(t, "one", d.Get(1))
	assert.Equal(t, "two", d.Get(2))

	var d2 container.Dictionary[TestKey, int]
	d2.Set(TestKey{Name: "alpha", Score: 100}, 1)
	d2.Set(TestKey{Name: "bravo", Score: 300}, 2)
	d2.Set(TestKey{Name: "charlie", Score: 700}, 8)
	assert.Equal(t, 1, d2.Get(TestKey{Name: "alpha", Score: 100}))
	assert.Equal(t, 2, d2.Get(TestKey{Name: "bravo", Score: 300}))
	assert.True(t, d2.Contains(TestKey{Name: "charlie", Score: 700}))
	assert.False(t, d2.Contains(TestKey{Name: "char_lie", Score: 700}))
	assert.False(t, d2.Contains(TestKey{Name: "charlie", Score: 100}))

	var d3 container.Dictionary[*TestKey, float64]
	k1 := TestKey{Name: "alpha", Score: 100}
	k2 := TestKey{Name: "bravo", Score: 300}
	k3 := TestKey{Name: "charlie", Score: 700}
	d3.Set(&k1, 1.1)
	d3.Set(&k2, 1.2)
	d3.Set(&k3, 1.3)
	assert.Equal(t, 1.1, d3.Get(&k1))
	assert.Equal(t, 1.2, d3.Get(&k2))
	assert.Equal(t, 1.3, d3.Get(&k3))
	assert.Equal(t, 0.0, d3.Get(nil))
	assert.False(t, d3.Contains(nil))
	assert.True(t, d3.Contains(&k3))
	d3.Remove(&k3)
	assert.False(t, d3.Contains(&k3))
	vf, ok := d3.Pop(&k2)
	assert.True(t, ok)
	assert.Equal(t, 1.2, vf)
	assert.False(t, d3.Contains(&k2))
	d3.Clear()
	assert.Equal(t, 0.0, d3.Get(&k1))
	assert.False(t, d3.Contains(&k1))
}
