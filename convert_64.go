package bitstring

// Uint32 returns the uint32 value represented by the 32 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint32(i uint) uint32 {
	bs.mustExist(i + 31)

	i64 := uint64(i)
	off := bitoffset(i64)
	loword := bs.data[wordoffset(i64)] >> off
	hiword := bs.data[wordoffset(i64+31)] & ((1 << off) - 1)
	return uint32(loword | hiword<<(uintsize-off))
}

// Uint64 returns the uint64 value represented by the 64 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Uint64(i uint) uint64 {
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

// SetUint32 sets the 32 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetUint32(i uint, x uint32) {
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
func (bs *Bitstring) SetUint64(i uint, x uint64) {
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
