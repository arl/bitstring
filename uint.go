package bitstring

import "fmt"

/* get */

// Uint8 returns the uint8 value represented by the 8 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint8(i int) uint8 {
	bs.mustExist(i + 7)

	return uint8(bs.uint(uint64(i), 7))
}

// Uint16 returns the uint16 value represented by the 16 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint16(i int) uint16 {
	bs.mustExist(i + 15)

	return uint16(bs.uint(uint64(i), 15))
}

// Uint32 returns the uint32 value represented by the 32 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint32(i int) uint32 {
	bs.mustExist(i + 31)

	return uint32(bs.uint(uint64(i), 31))
}

// Uint64 returns the uint64 value represented by the 64 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint64(i int) uint64 {
	bs.mustExist(i + 63)

	// fast path: i is a multiple of 64
	if i&((1<<6)-1) == 0 {
		return uint64(bs.data[i>>6])
	}

	i64 := uint64(i)
	w := wordoffset(i64)
	off := bitoffset(i64)
	loword := bs.data[w] >> off
	hiword := bs.data[w+1] & ((1 << off) - 1)
	return uint64(loword | hiword<<(uintsize-off))
}

func (bs *Bitstring) uint(i, nbits uint64) uint64 {
	off := bitoffset(i)
	loword := bs.data[wordoffset(i)] >> off
	hiword := bs.data[wordoffset(i+nbits)] & ((1 << off) - 1)
	return loword | hiword<<(uintsize-off)
}

// Uintn returns the n bits unsigned integer value represented by the n bits
// starting at the bit index i. It panics if there aren't enough bits in bs or
// if n is greater than the size of a machine word.
// TODO: reverse order of nbits and i params
func (bs *Bitstring) Uintn(n, i int) uint64 {
	if n > uintsize || n < 1 {
		panic(fmt.Sprintf("Uintn supports unsigned integers from 1 to %d bits long", uintsize))
	}
	bs.mustExist(i + n - 1)

	i64, n64 := uint64(i), uint64(n)
	j := wordoffset(i64)
	k := wordoffset(i64 + n64 - 1)
	looff := bitoffset(i64)
	loword := bs.data[j]
	if j == k {
		// fast path: same word
		return (loword >> looff) & lomask(n64)
	}
	hioff := bitoffset(i64 + n64)
	hiword := bs.data[k] & lomask(hioff)
	loword = himask(looff) & loword >> looff
	return loword | hiword<<(uintsize-looff)
}

/* set */

// SetUint8 sets the 8 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetUint8(i int, x uint8) {
	bs.mustExist(i + 7)

	i64 := uint64(i)
	lobit := bitoffset(i64)
	j := wordoffset(i64)
	k := wordoffset(i64 + 7)
	if j == k {
		// fast path: same word
		lobit := bitoffset(i64)
		neww := uint64(x) << lobit
		msk := mask(lobit, lobit+8)
		bs.data[j] = transferbits(bs.data[j], neww, msk)
		return
	}
	// transfer bits to low word
	bs.data[j] = transferbits(bs.data[j], uint64(x)<<lobit, himask(lobit))
	// transfer bits to high word
	lon := uintsize - lobit
	bs.data[k] = transferbits(bs.data[k], uint64(x)>>lon, lomask(8-lon))
}

// SetUint16 sets the 16 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetUint16(i int, x uint16) {
	bs.mustExist(i + 15)

	i64 := uint64(i)
	lobit := bitoffset(i64)
	j := wordoffset(i64)
	k := wordoffset(i64 + 15)
	if j == k {
		// fast path: same word
		neww := uint64(x) << lobit
		msk := mask(lobit, lobit+16)
		bs.data[j] = transferbits(bs.data[j], neww, msk)
		return
	}
	// transfer bits to low word
	bs.data[j] = transferbits(bs.data[j], uint64(x)<<lobit, himask(lobit))
	// transfer bits to high word
	lon := uintsize - lobit
	bs.data[k] = transferbits(bs.data[k], uint64(x)>>lon, lomask(16-lon))
}

// SetUint32 sets the 32 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetUint32(i int, x uint32) {
	bs.mustExist(i + 31)

	i64 := uint64(i)
	lobit := bitoffset(i64)
	j := wordoffset(i64)
	k := wordoffset(i64 + 31)
	if j == k {
		// fast path: same word
		neww := uint64(x) << lobit
		msk := mask(lobit, lobit+32)
		bs.data[j] = transferbits(bs.data[j], neww, msk)
		return
	}
	// transfer bits to low word
	bs.data[j] = transferbits(bs.data[j], uint64(x)<<lobit, himask(lobit))
	// transfer bits to high word
	lon := uintsize - lobit
	bs.data[k] = transferbits(bs.data[k], uint64(x)>>lon, lomask(32-lon))
}

// SetUint64 sets the 64 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetUint64(i int, x uint64) {
	bs.mustExist(i + 63)

	i64 := uint64(i)
	lobit := bitoffset(i64)
	j := wordoffset(i64)

	// fast path: i is a multiple of 64
	if i&((1<<6)-1) == 0 {
		bs.data[i>>6] = x
		return
	}

	k := wordoffset(i64 + 63)
	if j == k {
		// fast path: same word
		neww := uint64(x) << lobit
		msk := mask(lobit, lobit+64)
		bs.data[j] = transferbits(bs.data[j], neww, msk)
		return
	}
	// transfer bits to low word
	bs.data[j] = transferbits(bs.data[j], uint64(x)<<lobit, himask(lobit))
	// transfer bits to high word
	lon := (uintsize - lobit)
	bs.data[k] = transferbits(bs.data[k], uint64(x)>>lon, lomask(64-lon))
}

// SetUintn sets the n bits starting at i with the first n bits of value x.
// It panics if there aren't enough bits in bs or if n is greater than
// the size of a machine word.
func (bs *Bitstring) SetUintn(n, i int, x uint64) {
	if n > uintsize || n < 1 {
		panic(fmt.Sprintf("SetUintn supports unsigned integers from 1 to %d bits long", uintsize))
	}
	bs.mustExist(i + n - 1)

	i64, n64 := uint64(i), uint64(n)
	lobit := bitoffset(i64)
	j := wordoffset(i64)
	k := wordoffset(i64 + n64 - 1)
	if j == k {
		// fast path: same word
		x := (x & lomask(n64)) << lobit
		bs.data[j] = transferbits(bs.data[j], x, mask(lobit, lobit+n64))
		return
	}
	// slow path: first and last bits are on different words
	// transfer bits to low word
	lon := uintsize - lobit // how many bits of n we transfer to loword
	bs.data[j] = transferbits(bs.data[j], x<<lobit, himask(lon))

	// transfer bits to high word
	bs.data[k] = transferbits(bs.data[k], x>>lon, lomask(n64-lon))
}
