package container_test

import (
	"testing"

	"github.com/gabstv/container"
	"github.com/stretchr/testify/assert"
)

func TestDrain(t *testing.T) {
	ch := make(chan int, 50)
	defer close(ch)
	for i := 0; i < 50; i++ {
		ch <- i
	}
	items := container.Drain(ch)
	assert.Equal(t, 50, len(items))
	for i := 0; i < 50; i++ {
		assert.Equal(t, i, items[i])
	}
}
