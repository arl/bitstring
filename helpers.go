package bitstring

import (
	"fmt"
	"strings"
)

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
