// Package bitstring implements a fixed length bit string type and bit
// manipulation functions.

package bitstring

import (
	"fmt"
	"math/big"
	"math/bits"
	"math/rand"
)

// Bitstring implements a fixed-length bit string.
//
// Internally, bits are packed into an array of machine word integers. This
// implementation makes more efficient use of space than the alternative
// approach of using an array of booleans.
type Bitstring struct {
	// length in bits of the bit string
	length int

	// bits are packed in an array of uint64.
	data []uint64
}

// New creates a bit string of the specified length (in bits) with all bits
// initially set to zero (off).
func New(length int) *Bitstring {
	return &Bitstring{
		length: length,
		data:   make([]uint64, (length+64-1)/64),
	}
}

// Random creates a Bitstring of the length l in which each bit is assigned a
// random value using rng.
//
// Random randomly sets the uint32 values of the underlying slice, so it should
// be faster than creating a bit string and then randomly setting each
// individual bits.
func Random(length int, rng *rand.Rand) *Bitstring {
	bs := New(length)

	a := bs.data[:len(bs.data)] // remove bounds-checking

	// Fill words with random values
	for i := range a {
		a[i] = uint64(rng.Uint64())
	}

	// If the last word is not fully utilised, zero any out-of-bounds bits.
	// This is necessary because OnesCount and ZeroesCount count the
	// out-of-bounds bits.
	nused := bitoffset(uint64(length))
	if nused != 0 {
		mask := lomask(uint64(nused))
		a[len(a)-1] &= mask
	}
	return bs
}

// NewFromString returns the corresponding Bitstring for the given string of 1s
// and 0s in big endian order.
func NewFromString(s string) (*Bitstring, error) {
	bs := New(len(s))

	for i, c := range s {
		switch c {
		case '0':
			continue
		case '1':
			bs.SetBit(len(s) - i - 1)
		default:
			return nil, fmt.Errorf("illegal character at position %v: %#U", i, c)
		}
	}
	return bs, nil
}

// Len returns the length if bs, that is the number of bits it contains.
func (bs *Bitstring) Len() int {
	return int(bs.length)
}

// Data returns the bitstring underlying slice.
func (bs *Bitstring) Data() []uint64 {
	return bs.data
}

// Bit returns a boolean indicating whether the bit at index i is set or not.
//
// If i is greater than the bitstring length, Bit will panic.
func (bs *Bitstring) Bit(i int) bool {
	bs.mustExist(i)

	w := wordoffset(uint64(i))
	off := bitoffset(uint64(i))
	mask := bitmask(off)
	return (bs.data[w] & mask) != 0
}

// SetBit sets the bit at index i.
//
// If i is greater than the bitstring length, SetBit will panic.
func (bs *Bitstring) SetBit(i int) {
	bs.mustExist(i)

	w := wordoffset(uint64(i))
	off := bitoffset(uint64(i))
	bs.data[w] |= bitmask(off)
}

// ClearBit clears the bit at index i.
//
// If i is greater than the bitstring length, ClearBit will panic.
func (bs *Bitstring) ClearBit(i int) {
	bs.mustExist(i)

	w := wordoffset(uint64(i))
	off := bitoffset(uint64(i))
	bs.data[w] &= ^bitmask(off)
}

// FlipBit flips (i.e toggles) the bit at index i.
//
// If i is greater than the bitstring length, FlipBit will panic.
func (bs *Bitstring) FlipBit(i int) {
	bs.mustExist(i)

	w := wordoffset(uint64(i))
	off := bitoffset(uint64(i))
	bs.data[w] ^= (1 << off)
}

// OnesCount counts the number of one bits.
func (bs *Bitstring) OnesCount() int {
	var count int
	for _, x := range bs.data {
		count += bits.OnesCount64(x)
	}
	return count
}

// ZeroesCount counts the number of zero bits.
func (bs *Bitstring) ZeroesCount() int {
	return bs.length - bs.OnesCount()
}

func alignedRev(data []uint64) {
	p := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&data)).Data)

	var buf []byte
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	hdr.Data = uintptr(p)
	hdr.Len = len(data) * 8
	hdr.Cap = len(data) * 8

	// NOTE: we don't care about native endianness here since we're performing a
	// byte-per-byte swap.
	for i := 0; i < len(buf)/2; i++ {
		buf[i], buf[len(buf)-i-1] = reverseLut[buf[len(buf)-i-1]], reverseLut[buf[i]]
	}
}

// Reverse reverses all bits in-place.
func (bs *Bitstring) Reverse() {
	alignedRev(bs.data)

	if bs.length%64 == 0 {
		// Fast path, word-aligned reversal.
		return
	}

	off := bitoffset(uint64(bs.length))
	shift := 64 - off
	savemask := lomask(shift)
	mask := uint64(0)

	for i := len(bs.data) - 1; i >= 0; i-- {
		save := (bs.data[i] & savemask) << off
		bs.data[i] >>= shift
		bs.data[i] |= mask
		mask = save
	}
}
// BigInt returns the big.Int representation of bs.
func (bs *Bitstring) BigInt() *big.Int {
	bi := new(big.Int)
	_, _ = bi.SetString(bs.String(), 2)

	return bi
}

// String returns a string representation of bs in big endian order.
func (bs *Bitstring) String() string {
	b := make([]byte, bs.length)
	for i := 0; i < bs.length; i++ {
		if bs.Bit(i) {
			b[bs.length-1-i] = '1'
		} else {
			b[bs.length-1-i] = '0'
		}
	}
	return string(b)
}

// Clone creates and returns a new Bitstring that is a clone of src.
func Clone(src *Bitstring) *Bitstring {
	dst := make([]uint64, len(src.data))
	copy(dst, src.data)
	return &Bitstring{
		length: src.length,
		data:   dst,
	}
}

// Copy copies a source Bitstring into a destination Bitstring, shrinking or
// expanding it if necessary.
func Copy(dst, src *Bitstring) {
	switch {
	case dst.length == src.length:
	case dst.length < src.length:
		// XXX: Reallocate the whole bitstring, but is it really faster?
		dst.data = make([]uint64, len(src.data))
		dst.length = src.length
	case dst.length > src.length:
		dst.data = dst.data[:len(src.data)]
		dst.length = src.length
	}

	copy(dst.data, src.data)
}

// Equals returns true if bs and other have the same length and each bit are
// identical, or if bs and other both point to the same Bitstring instance (i.e
// pointer equality).
func (bs *Bitstring) Equals(other *Bitstring) bool {
	switch {
	case bs == other:
		return true
	case bs != nil && other == nil,
		bs == nil && other != nil:
		return false
	case bs.length == other.length:
		return u64cmp(bs.data, other.data)
	}
	return false
}

/*
// RotateLeft rotates the bitstring by (k mod len) bits.
func (bs *Bitstring) RotateLeft(k int) {
	panic("unimplemented")
}

// RotateRight rotates the bitstring by (k mod len) bits.
func (bs *Bitstring) RotateRight(k int) {
	panic("unimplemented")
}

// Reverse reverses the order of bits in the bitstring.
func (bs *Bitstring) Reverse(k int) {
	panic("unimplemented")
}

// LeadingZeros returns the number of leading zero bits in the bitstring.
func (bs *Bitstring) LeadingZeros() int {
	panic("unimplemented")
}

// TrailingZeros returns the number of trailing zero bits in the bitstring.
func (bs *Bitstring) TrailingZeros() int {
	panic("unimplemented")
}
*/
