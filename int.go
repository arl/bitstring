package bitstring

/* unsigned integer get */

// Uint8 interprets the 8 bits at offset off as an uint8 in big endian and
// returns its value. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) Uint8(off int) uint8 {
	bs.mustExist(off + 7)

	return uint8(bs.uint(uint64(off), 7))
}

// Uint16 interprets the 16 bits at offset off as an uint16 in big endian and
// returns its value. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) Uint16(off int) uint16 {
	bs.mustExist(off + 15)

	return uint16(bs.uint(uint64(off), 15))
}

// Uint32 interprets the 32 bits at offset off as an uint32 in big endian and
// returns its value. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) Uint32(off int) uint32 {
	bs.mustExist(off + 31)

	return uint32(bs.uint(uint64(off), 31))
}

// Uint64 interprets the 64 bits at offset off as an uint64 in big endian and
// returns its value. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) Uint64(off int) uint64 {
	bs.mustExist(off + 63)

	if off&((1<<6)-1) == 0 {
		// Fast path: off is a multiple of 64.
		return uint64(bs.data[off>>6])
	}

	i := uint64(off)
	w := wordoffset(i)
	bit := bitoffset(i)
	loword := bs.data[w] >> bit
	hiword := bs.data[w+1] & ((1 << bit) - 1)
	return uint64(loword | hiword<<(64-bit))
}

// Uintn interprets the n bits at offset off as an n-bit unsigned integer in big
// endian and returns its value. Behavior is undefined if there aren't enough
// bits. Panics if nbits is greater than 64.
func (bs *Bitstring) Uintn(off, n int) uint64 {
	if n > 64 || n < 1 {
		panic("Uintn supports unsigned integers from 1 to 64 bits long")
	}
	bs.mustExist(off + n - 1)

	i, nbits := uint64(off), uint64(n)
	j := wordoffset(i)
	k := wordoffset(i + nbits - 1)
	looff := bitoffset(i)
	loword := bs.data[j]
	if j == k {
		// Fast path: value doesn't cross uint64 boundaries.
		return (loword >> looff) & lomask(nbits)
	}
	hioff := bitoffset(i + nbits)
	hiword := bs.data[k] & lomask(hioff)
	loword = himask(looff) & loword >> looff
	return loword | hiword<<(64-looff)
}

func (bs *Bitstring) uint(off, n uint64) uint64 {
	bit := bitoffset(off)
	loword := bs.data[wordoffset(off)] >> bit
	hiword := bs.data[wordoffset(off+n)] & ((1 << bit) - 1)
	return loword | hiword<<(64-bit)
}

/* unsigned integer set */

// SetUint8 sets the 8 bits at offset off with the given int8 value, in big
// endian. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) SetUint8(off int, val uint8) {
	bs.mustExist(off + 7)

	i := uint64(off)
	lobit := bitoffset(i)
	j := wordoffset(i)
	k := wordoffset(i + 7)
	if j == k {
		// Fast path: value doesn't cross uint64 boundaries.
		lobit := bitoffset(i)
		neww := uint64(val) << lobit
		msk := mask(lobit, lobit+8)
		bs.data[j] = transferbits(bs.data[j], neww, msk)
		return
	}
	// Transfer bits to low word.
	bs.data[j] = transferbits(bs.data[j], uint64(val)<<lobit, himask(lobit))
	// Transfer bits to high word.
	lon := 64 - lobit
	bs.data[k] = transferbits(bs.data[k], uint64(val)>>lon, lomask(8-lon))
}

// SetUint16 sets the 8 bits at offset off with the given int8 value, in big
// endian. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) SetUint16(off int, val uint16) {
	bs.mustExist(off + 15)

	i := uint64(off)
	lobit := bitoffset(i)
	j := wordoffset(i)
	k := wordoffset(i + 15)
	if j == k {
		// Fast path: value doesn't cross uint64 boundaries.
		neww := uint64(val) << lobit
		msk := mask(lobit, lobit+16)
		bs.data[j] = transferbits(bs.data[j], neww, msk)
		return
	}
	// Transfer bits to low word.
	bs.data[j] = transferbits(bs.data[j], uint64(val)<<lobit, himask(lobit))
	// Transfer bits to high word.
	lon := 64 - lobit
	bs.data[k] = transferbits(bs.data[k], uint64(val)>>lon, lomask(16-lon))
}

// SetUint32 sets the 8 bits at offset off with the given int8 value, in big
// endian. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) SetUint32(off int, val uint32) {
	bs.mustExist(off + 31)

	i := uint64(off)
	lobit := bitoffset(i)
	j := wordoffset(i)
	k := wordoffset(i + 31)
	if j == k {
		// Fast path: value doesn't cross uint64 boundaries.
		neww := uint64(val) << lobit
		msk := mask(lobit, lobit+32)
		bs.data[j] = transferbits(bs.data[j], neww, msk)
		return
	}
	// Transfer bits to low word.
	bs.data[j] = transferbits(bs.data[j], uint64(val)<<lobit, himask(lobit))
	// Transfer bits to high word.
	lon := 64 - lobit
	bs.data[k] = transferbits(bs.data[k], uint64(val)>>lon, lomask(32-lon))
}

// SetUint64 sets the 8 bits at offset off with the given int8 value, in big
// endian. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) SetUint64(off int, val uint64) {
	bs.mustExist(off + 63)

	i := uint64(off)
	lobit := bitoffset(i)
	j := wordoffset(i)

	if off&((1<<6)-1) == 0 {
		// Fast path: off is a multiple of 64.
		bs.data[off>>6] = val
		return
	}

	// Transfer bits to low word.
	bs.data[j] = transferbits(bs.data[j], uint64(val)<<lobit, himask(lobit))
	// Transfer bits to high word.
	lon := (64 - lobit)
	k := wordoffset(i + 63)
	bs.data[k] = transferbits(bs.data[k], uint64(val)>>lon, lomask(64-lon))
}

// SetUintn sets the n bits at offset off with the given n-bit unsigned integer in
// big endian. Behavior is undefined if there aren't enough bits. Panics if
// nbits is greater than 64.
func (bs *Bitstring) SetUintn(off, n int, val uint64) {
	if n > 64 || n < 1 {
		panic("SetUintn supports unsigned integers from 1 to 64 bits long")
	}
	bs.mustExist(off + n - 1)

	i, nbits := uint64(off), uint64(n)
	lobit := bitoffset(i)
	j := wordoffset(i)
	k := wordoffset(i + nbits - 1)
	if j == k {
		// Fast path: value doesn't cross uint64 boundaries.
		x := (val & lomask(nbits)) << lobit
		bs.data[j] = transferbits(bs.data[j], x, mask(lobit, lobit+nbits))
		return
	}

	// First and last bits are on different words.
	// Transfer bits to low word.
	lon := 64 - lobit // how many bits of n we transfer to loword
	bs.data[j] = transferbits(bs.data[j], val<<lobit, himask(lon))

	// Transfer bits to high word.
	bs.data[k] = transferbits(bs.data[k], val>>lon, lomask(nbits-lon))
}

/* signed get */

// Int8 interprets the 8 bits at offset off as an int8 in big endian and
// returns its value. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) Int8(off int) int8 { return int8(bs.Uint8(off)) }

// Int16 interprets the 16 bits at offset off as an int16 in big endian and
// returns its value. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) Int16(off int) int16 { return int16(bs.Uint16(off)) }

// Int32 interprets the 32 bits at offset off as an int32 in big endian and
// returns its value. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) Int32(off int) int32 { return int32(bs.Uint32(off)) }

// Int64 interprets the 64 bits at offset off as an int64 in big endian and
// returns its value. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) Int64(off int) int64 { return int64(bs.Uint64(off)) }

// Intn interprets the n bits at offset off as an n-bit signed integer in big
// endian and returns its value. Behavior is undefined if there aren't enough
// bits. Panics if nbits is greater than 64.
func (bs *Bitstring) Intn(off, n int) int64 { return int64(bs.Uintn(off, n)) }

/* signed integer set */

// SetInt8 sets the 8 bits at offset off with the given int8 value, in big
// endian. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) SetInt8(off int, val int8) { bs.SetUint8(off, uint8(val)) }

// SetInt16 sets the 16 bits at offset off with the given int16 value, in big
// endian. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) SetInt16(off int, val int16) { bs.SetUint16(off, uint16(val)) }

// SetInt32 sets the 32 bits at offset off with the given int32 value, in big
// endian. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) SetInt32(off int, val int32) { bs.SetUint32(off, uint32(val)) }

// SetInt64 sets the 64 bits at offset off with the given int64 value, in big
// endian. Behavior is undefined if there aren't enough bits.
func (bs *Bitstring) SetInt64(off int, val int64) { bs.SetUint64(off, uint64(val)) }

// SetIntn sets the n bits at offset off with the given n-bit signed integer in
// big endian. Behavior is undefined if there aren't enough bits. Panics if
// nbits is greater than 64.
func (bs *Bitstring) SetIntn(off, n int, val int64) { bs.SetUintn(off, n, uint64(val)) }
