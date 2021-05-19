// +build bitstring_debug

package bitstring

import "fmt"

// mustExist panics if i is not a valid bit index for bs, that is if i is
// greater than bs.length.
func (bs *Bitstring) mustExist(i int) {
	if i >= bs.length || i < 0 {
		panic(fmt.Sprintf("Bitstring: index %d is out of range [%d, %d]", i, 0, bs.length))
	}
}
