// Package bitstring implements a fixed length bit string type and bit
// manipulation functions.

package bitstring

import (
	"fmt"
	"math/big"
	"math/bits"
	"math/rand"
	"reflect"
	"unsafe"
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

// Len returns the length of bs in bits.
func (bs *Bitstring) Len() int {
	return int(bs.length)
}

// Bits returns raw access to the bitstring underlying uint64 slice. The result
// and bs share the same underlying array.
//
// Bits is intended to support implementation of missing low-level Bitstring
// functionality outside this package; it should be avoided otherwise.
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
	// We first reverse the whole bitstring.
	alignedRev(bs.data)

	if bs.length%64 == 0 {
		// No need for extra shifting.
		return
	}

	rightShiftBits(bs.data, bitoffset(uint64(64-bs.length)))
}

// Flip flips all bits (replaces ones with zeroes and zeroes with ones).
func (bs *Bitstring) Flip() {
	bs.FlipRange(0, bs.length)
}

// NewFromBig creates a new Bitstring using the absolute value of the big.Int
// bi.
//
// The number of bits of the new Bitstring depends on the number of significant
// bits in the binary representation of bi.
func NewFromBig(bi *big.Int) *Bitstring {
	words := bi.Bits()

	p := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&words)).Data)

	var bigData []uint64
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&bigData))
	hdr.Data = uintptr(p)
	hdr.Len = len(words)
	hdr.Cap = hdr.Len

	datalen := bi.BitLen() / 64
	if bi.BitLen()%64 != 0 {
		datalen++
	}

	bs := &Bitstring{
		length: bi.BitLen(),
		data:   make([]uint64, datalen),
	}

	copy(bs.data, bigData)

	return bs
}

// BigInt returns the big.Int representation of bs.
func (bs *Bitstring) BigInt() *big.Int {
	cpy := bs.Clone()

	p := unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&cpy.data)).Data)

	var words []big.Word
	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&words))
	hdr.Data = uintptr(p)
	hdr.Len = len(cpy.data) * (64 / wordsize)
	hdr.Cap = hdr.Len

	bint := new(big.Int)
	bint.SetBits(words)

	return bint
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
func (bs *Bitstring) Clone() *Bitstring {
	dst := make([]uint64, len(bs.data))
	copy(dst, bs.data)
	return &Bitstring{
		length: bs.length,
		data:   dst,
	}
}

// Copy copies a source Bitstring into a destination Bitstring, shrinking or
// expanding it if necessary.
func Copy(dst, src *Bitstring) {
	switch {
	case dst.length == src.length:
	case dst.length < src.length:
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

// LeadingZeroes returns the number of leading bits that are set to 0 in bs.
// (i.e the number of consecutives 0's starting from the MSB (most significant
// bit).
func (bs *Bitstring) LeadingZeroes() int {
	bitoff := int(bitoffset(uint64(bs.length)))
	start := len(bs.data) - 1

	n := 0
	for i := start; i >= 0; i-- {
		leading := bits.LeadingZeros64(bs.data[i])

		// We treat the first word separately if the bistring length is not a
		// multiple of the wordsize because in this case we must omit the 'extra
		// bits' from the count the count of leading zeroes.
		if i == start && bitoff != 0 {
			// Limit to 'off' the number of bits we count.
			leading -= 64 - bitoff
			n += leading
			if leading != bitoff {
				break // early exit if useful bits are not all 0s.
			}
		} else {
			// Subsequent words
			n += leading
			if leading != 64 {
				break
			}
		}
	}

	return n
}

// LeadingOnes returns the number of leading bits that are set to 1 in bs. (i.e
// the number of consecutives 1's starting from the MSB (most significant bit).
func (bs *Bitstring) LeadingOnes() int {
	bitoff := int(bitoffset(uint64(bs.length)))
	start := len(bs.data) - 1

	n := 0
	for i := start; i >= 0; i-- {
		// We treat the first word separately if the bistring length is not a
		// multiple of the wordsize because in this case we must omit the 'extra
		// bits' from the count the count of leading zeroes.
		if i == start && bitoff != 0 {
			leading := bits.LeadingZeros64(lomask(uint64(bitoff)) ^ bs.data[i])

			// Limit to 'off' the number of bits we count.
			leading -= 64 - bitoff
			n += leading
			if leading != bitoff {
				break // early exit if useful bits are not all 0s.
			}
		} else {
			leading := bits.LeadingZeros64(^bs.data[i])

			// Subsequent words
			n += leading
			if leading != 64 {
				break
			}
		}
	}

	return n
}

// TrailingZeroes returns the number of trailing bits that are set to 0 in bs.
// (i.e the number of consecutives 0's starting from the LSB (least significant
// bit).
func (bs *Bitstring) TrailingZeroes() int {
	bitoff := int(bitoffset(uint64(bs.length)))
	last := len(bs.data) - 1

	n := 0
	for i := 0; i < len(bs.data); i++ {
		trailing := bits.TrailingZeros64(bs.data[i])

		if i == last && bitoff != 0 && trailing == 64 {
			// There's one specific case we need to take care of: if the last
			// word is 0 and the bitstring length is not a multiple of 64 then
			// the actual number of trailing bits is not 64, we need to limit it
			// to the number of useful bits only.
			trailing = bitoff
		}

		n += trailing
		if trailing != 64 {
			break
		}
	}
	return n
}

// TrailingOnes returns the number of trailing bits that are set to 1 in bs.
// (i.e the number of consecutives 1's starting from the LSB (least significant
// bit).
func (bs *Bitstring) TrailingOnes() int {
	bitoff := int(bitoffset(uint64(bs.length)))
	last := len(bs.data) - 1

	n := 0
	for i := 0; i < len(bs.data); i++ {
		trailing := bits.TrailingZeros64(^bs.data[i])

		if i == last && bitoff != 0 && trailing == 64 {
			// There's one specific case we need to take care of: if the last
			// word if 0 and the bitstring length is not a multiple of the
			// wordsize, then the effective number of trailing bits is not 64,
			// we need to limit it to the number of useful bits only.
			trailing = bitoff
		}

		n += trailing
		if trailing != 64 {
			break
		}
	}
	return n
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
*/
