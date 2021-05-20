// +build bitstring_debug

package bitstring

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetBitDebug(t *testing.T) {
	bs := New(1)
	t.Run("panics on index too high", func(t *testing.T) {
		// The index of an individual bit must be within the range 0 to
		// length-1.
		assert.Panics(t, func() { bs.ClearBit(1) })
	})
}
