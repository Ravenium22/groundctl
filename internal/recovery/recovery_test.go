package recovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrapNoPanic(t *testing.T) {
	// Should not panic
	func() {
		defer Wrap("test")()
		// normal execution
	}()
}

func TestWrapFunctionReturned(t *testing.T) {
	fn := Wrap("test")
	assert.NotNil(t, fn)
}
