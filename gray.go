package bitstring

// Gray8 returns the uint8 value represented by the 8 gray-coded bits starting
// at the given bit. It panics if there are not enough bits.
func (bs *Bitstring) Gray8(i int) uint8 {
	v := bs.Uint8(i)
	v ^= v >> 4
	v ^= v >> 2
	v ^= v >> 1
	return v
}

// Gray16 returns the uint8 value represented by the 16 gray-coded bits starting
// at the given bit. It panics if there are not enough bits.
func (bs *Bitstring) Gray16(i int) uint16 {
	v := bs.Uint16(i)
	v ^= v >> 8
	v ^= v >> 4
	v ^= v >> 2
	v ^= v >> 1
	return v
}

// Gray32 returns the uint32 value represented by the 32 gray-coded bits starting
// at the given bit. It panics if there are not enough bits.
func (bs *Bitstring) Gray32(i int) uint32 {
	v := bs.Uint32(i)
	v ^= v >> 16
	v ^= v >> 8
	v ^= v >> 4
	v ^= v >> 2
	v ^= v >> 1
	return v
}

// Gray64 returns the uint64 value represented by the 64 gray-coded bits starting
// at the given bit. It panics if there are not enough bits.
func (bs *Bitstring) Gray64(i int) uint64 {
	v := bs.Uint64(i)
	v ^= v >> 32
	v ^= v >> 16
	v ^= v >> 8
	v ^= v >> 4
	v ^= v >> 2
	v ^= v >> 1
	return v
}

// Grayn returns the n-bit unsigned integer value represented by the n
// gray-coded bits starting at the bit index i. It panics if there are not
// enough bits or if n is greater than the size of a machine word.
func (bs *Bitstring) Grayn(nbits, i int) uint64 {
	v := bs.Uintn(nbits, i)
	v ^= v >> 32
	v ^= v >> 16
	v ^= v >> 8
	v ^= v >> 4
	v ^= v >> 2
	v ^= v >> 1
	return v
}
