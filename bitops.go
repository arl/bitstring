package bitstring

import (
	"math"
	"unsafe"
)

const wordsize = 32 << (^uint(0) >> 63) // 32 or 64

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
// behavior if n is greater than 64.
func lomask(n uint64) uint64 { return math.MaxUint64 >> (64 - n) }

// himask returns a mask to keep the n MSB (most significant bits). Undefined
// behavior if n is greater than 64.
func himask(n uint64) uint64 { return math.MaxUint64 << n }

// transferbits returns the word that results from transferring some bits from
// src to dst, where set bits in mask specify the bits to transfer.
func transferbits(dst, src, mask uint64) uint64 {
	return dst&^mask | src&mask
}

// lsb returns the offset of the lowest significant set bit in v. That is, the
// index of the rightmost 1.
//
// Note: lsb(0) is meaningless, it's the caller responsibility to not use the
// result of lsb(0).
func lsb(v uint64) uint64 {
	var num uint64

	if (v & 0xffffffff) == 0 {
		num += 32
		v >>= 32
	}
	if (v & 0xffff) == 0 {
		num += 16
		v >>= 16
	}
	if (v & 0xff) == 0 {
		num += 8
		v >>= 8
	}
	if (v & 0xf) == 0 {
		num += 4
		v >>= 4
	}
	if (v & 0x3) == 0 {
		num += 2
		v >>= 2
	}
	if (v & 0x1) == 0 {
		num++
	}
	return num
}

// msb returns the offset of the most significant set bit in v. That is, the
// index of the leftmost 1.
//
// Note: msb(0) is meaningless, it's the caller responsibility to not use the
// result of msb(0).
func msb(v uint64) uint64 {
	var num uint64

	if (v & 0xffffffff00000000) != 0 {
		num += 32
		v >>= 32
	}
	if (v & 0xffff0000) != 0 {
		num += 16
		v >>= 16
	}
	if (v & 0xff00) != 0 {
		num += 8
		v >>= 8
	}
	if (v & 0xf0) != 0 {
		num += 4
		v >>= 4
	}
	if (v & 0xc) != 0 {
		num += 2
		v >>= 2
	}
	if (v & 0x2) != 0 {
		num++
		v >>= 1
	}

	return num
}

// fastmsbLittleEndian is faster version of msb that only works on little endian
// architectures. About 50% faster than msb on amd64. Rely on the fact that Go
// uses IEEE 754 floating point representation. Converts v to float64, then
// extracts the exponent bits of the IEEE754 representation.
func fastmsbLittleEndian(v uint64) uint64 {
	if v == math.MaxUint64 {
		return 63
	}

	f := float64(v)
	arr := *(*[2]uint32)(unsafe.Pointer(&f))
	return uint64(arr[1]>>20 - 1023)
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
