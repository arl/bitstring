package bitstring

import (
	"fmt"
	"strconv"
	"testing"
)

func atobin(s string) uint64 {
	i, err := strconv.ParseUint(s, 2, 64)
	if err != nil {
		panic(fmt.Sprintf("Can't convert %s to base 2: %s", s, err))
	}
	return i
}

func equalbits(tb testing.TB, got, want *Bitstring) {
	tb.Helper()

	if want.length != got.length {
		tb.Fatalf("got bitstring with length = %d, want length = %d", got.length, want.length)
	}

	for i := 0; i < got.Len(); i++ {
		if want.Bit(i) != got.Bit(i) {
			tb.Fatalf("bitstrings differs from expected from bit %d (bit %d of word %d)", i, bitoffset(uint64(i)), wordoffset(uint64(i)))
		}
	}

	// Even if all useful bits are equal we want the extra bits to be the same
	// (0's) of that could mess up with functions such as OnesCount, etc.
	for i := range want.data {
		if want.data[i] != got.data[i] {
			tb.Fatalf("differences in extra bits i=%d got=0v%b want=0b%b", i, got.data[i], want.data[i])
		}
	}
}
