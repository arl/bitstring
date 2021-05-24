package bitstring

import (
	"encoding/binary"
	"fmt"
	"strings"
	"unsafe"
)

// bitmask returns a mask where only the nth bit of a word is set.
func bitmask(n uint64) uint64 { return 1 << n }

// wordoffset returns, for a given bit n of a bit string, the offset
// of the word that contains bit n.
func wordoffset(n uint64) uint64 { return n / 64 }

// bitoffset returns, for a given bit n of a bit string, the offset of that bit
// with respect to the first bit of the word that contains it.
func bitoffset(n uint64) uint64 { return n & (64 - 1) }

// mask returns a mask that keeps the bits in the range [l, h) behavior
// undefined if any argument is greater than the size of a machine word.
func mask(l, h uint64) uint64 { return lomask(h) & himask(l) }

// lomask returns a mask to keep the n LSB (least significant bits). Undefined
// behavior if n is greater than uintsize.
func lomask(n uint64) uint64 { return maxuint >> (64 - n) }

// himask returns a mask to keep the n MSB (most significant bits). Undefined
// behavior if n is greater than uintsize.
func himask(n uint64) uint64 { return maxuint << n }

// transferbits returns the word that results from transferring some bits from
// src to dst, where set bits in mask specify the bits to transfer.
func transferbits(dst, src, mask uint64) uint64 {
	return dst&^mask | src&mask
}

// firstSetBit returns the offset of the first set bit in w
func firstSetBit(w uint64) uint64 {
	var num uint64

	if (w & 0xffffffff) == 0 {
		num += 32
		w >>= 32
	}
	if (w & 0xffff) == 0 {
		num += 16
		w >>= 16
	}
	if (w & 0xff) == 0 {
		num += 8
		w >>= 8
	}
	if (w & 0xf) == 0 {
		num += 4
		w >>= 4
	}
	if (w & 0x3) == 0 {
		num += 2
		w >>= 2
	}
	if (w & 0x1) == 0 {
		num++
	}
	return num
}

func reverseBytes(buf []byte) []byte {
	for i := 0; i < len(buf)/2; i++ {
		buf[i], buf[len(buf)-i-1] = buf[len(buf)-i-1], buf[i]
	}
	return buf
}

// shifts all words of n bits to the right. invariant: 0 <= off <= 64.
func rightShiftBits(words []uint64, n uint64) {
	shift := 64 - n
	mask := lomask(n)
	prev := uint64(0) // bits from previous word

	for i := len(words) - 1; i >= 0; i-- {
		save := (words[i] & mask) << shift
		words[i] >>= n
		words[i] |= prev
		prev = save
	}
}

// if n is a power of 2, ispow2 returns (v, true) such that (1<<v) gives n, or
// (0, false) if n is not a power of 2.
//
// panics if n == 0
func ispow2(n uint64) (uint64, bool) {
	if (n & -n) != n {
		// n is not a power of 2
		return 0, false
	}

	for i := uint64(0); i < 64; i++ {
		if n == 1<<i {
			return i, true
		}
	}

	panic("unreachable")
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

var nativeEndian binary.ByteOrder

func init() {
	i := uint32(1)
	b := (*[4]byte)(unsafe.Pointer(&i))
	if b[0] == 1 {
		nativeEndian = binary.LittleEndian
	} else {
		nativeEndian = binary.BigEndian
	}
}
