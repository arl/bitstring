// Package bitstring implements a fixed length bit string type and bit string
// manipulation functions
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

	// fill words with random values
	switch uintsize {
	case 32:
		for i := range a {
			a[i] = uint64(rng.Uint32())
		}
	case 64:
		for i := range a {
			a[i] = uint64(rng.Uint64())
		}
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

// MakeFromString returns the corresponding Bitstring for the given string of 1s
// and 0s in big endian order.
func MakeFromString(s string) (*Bitstring, error) {
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

// Len returns the lenght if bs, that is the number of bits it contains.
func (bs *Bitstring) Len() int {
	return int(bs.length)
}

// Data returns the bitstring underlying slice.
func (bs *Bitstring) Data() []uint64 {
	return bs.data
}

// Bit returns a boolean indicating wether the bit at index i is set or not.
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

// BigInt returns the big.Int representation of bs.
func (bs *Bitstring) BigInt() *big.Int {
	bi := new(big.Int)
	if _, ok := bi.SetString(bs.String(), 2); !ok {
		panic(fmt.Sprintf("couldn't convert bit string \"%s\" to big.Int", bs.String()))
	}
	return bi
}

func minuint(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

// SwapRange swaps the same range of bits between bs and other.
//
// Both Bitstring may have different length the bit with index start+lenght must
// be valid on both or SwapRange will panic.
func SwapRange(bs1, bs2 *Bitstring, start, length int) {
	bs1.mustExist(start + length - 1)
	bs2.mustExist(start + length - 1)

	// swap the required bits of the first word
	start64, len64 := uint64(start), uint64(length)
	i := start64 / uintsize
	start64 = bitoffset(start64)
	end := minuint(start64+len64, uintsize)
	remain := len64 - (end - start64)
	swapBits(bs1, bs2, i, mask(start64, end))
	i++

	// swap whole words but the last one
	for remain > uintsize {
		bs1.data[i], bs2.data[i] = bs2.data[i], bs1.data[i]
		remain -= uintsize
		i++
	}

	// swap the remaining bits of the last word
	if remain != 0 {
		swapBits(bs1, bs2, i, lomask(remain))
	}
}

// swapBits swaps range of bits from one word to another.
// w is the index of the word containing the bits to swap, and m is a mask that specifies
// whilch bits of that word will be swapped.
func swapBits(x, y *Bitstring, w, mask uint64) {
	keep := ^mask
	xkeep, ykeep := x.data[w]&keep, y.data[w]&keep
	xswap, yswap := x.data[w]&mask, y.data[w]&mask
	x.data[w] = xkeep | yswap
	y.data[w] = ykeep | xswap
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

// Copy creates and returns a new Bitstring that is a copy of src.
//
// TODO: rename Clone and make Copy a free functions like copy builtin
// (bitstring.Copy(dst, src)) with fewer allocation possible (possibly 0).
func Copy(src *Bitstring) *Bitstring {
	dst := make([]uint64, len(src.data))
	copy(dst, src.data)
	return &Bitstring{
		length: src.length,
		data:   dst,
	}
}

// Equals returns true if bs and other have the same lenght and each bit are
// identical, or if bs and other both point to the same Bitstring instance
// (i.e pointer equality).
func (bs *Bitstring) Equals(other *Bitstring) bool {
	switch {
	case bs == other:
		return true
	case bs != nil && other == nil,
		bs == nil && other != nil:
		return false
	case bs.length == other.length:
		od := other.data[:len(bs.data)] // remove BCE
		for i, v := range bs.data {
			if v != od[i] {
				return false
			}
		}
		return true
	}
	return false
}
