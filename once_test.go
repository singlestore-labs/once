package once

import (
	"testing"

	"gotest.tools/assert"
)

func TestOnce(t *testing.T) {
	var counter int
	once := New(func() { counter++ })
	once.Do()
	once.Do()
	once.Do()
	once.Do()
	once.Do()

	assert.Equal(t, 1, counter)
}
