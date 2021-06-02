package bitstring

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"unsafe"
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

func sprintbuf(b []byte) string {
	var sb strings.Builder
	for i := range b {
		fmt.Fprintf(&sb, "%08b ", b[i])
	}

	return sb.String()
}

func printbuf(b []byte) {
	fmt.Println(sprintbuf(b))
}

// returns a string representing the first n bits of the base-2 representation
// of val (unsigned).
func sprintubits(val uint64, nbits int) string {
	return fmt.Sprintf(fmt.Sprintf("%%0%db", nbits), val)
}

// returns a string representing the first n bits of the base-2 representation
// of val (signed).
func sprintsbits(val int64, nbits int) string {
	if val < 0 {
		// casting to uint will show us the 2's complement
		return sprintubits(uint64(val), nbits)
	}
	return fmt.Sprintf(fmt.Sprintf("%%0%db", nbits), val)
}

// prints a string representing the first n bits of the base-2 representation of val.
//lint:ignore U1000 (unused but useful for debugging)
func printbits(val, n int) {
	fmt.Printf(fmt.Sprintf("%%0%db\n", n), val)
}

var nativeEndian binary.ByteOrder

func init() {
	i := uint32(1)
	b := (*[4]byte)(unsafe.Pointer(&i))
	if b[0] == 1 {
		nativeEndian = binary.LittleEndian
	} else {
		nativeEndian = binary.BigEndian
	}
	fmt.Println(nativeEndian, wordsize)
}
