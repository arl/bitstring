package bitstring

// Gray8 interprets the 8 bits at offset off as a gray-coded uint8 in big
// endian and returns its value. Behavior is undefined if there aren't enough
// bits.
func (bs *Bitstring) Gray8(val int) uint8 {
	v := bs.Uint8(val)
	v ^= v >> 4
	v ^= v >> 2
	v ^= v >> 1
	return v
}

// Gray16 interprets the 16 bits at offset off as a gray-coded uint16 in big
// endian and returns its value. Behavior is undefined if there aren't enough
// bits.
func (bs *Bitstring) Gray16(val int) uint16 {
	v := bs.Uint16(val)
	v ^= v >> 8
	v ^= v >> 4
	v ^= v >> 2
	v ^= v >> 1
	return v
}

// Gray32 interprets the 32 bits at offset off as a gray-coded uint32 in big
// endian and returns its value. Behavior is undefined if there aren't enough
// bits.
func (bs *Bitstring) Gray32(val int) uint32 {
	v := bs.Uint32(val)
	v ^= v >> 16
	v ^= v >> 8
	v ^= v >> 4
	v ^= v >> 2
	v ^= v >> 1
	return v
}

// Gray64 interprets the 64 bits at offset off as a gray-coded uint64 in big
// endian and returns its value. Behavior is undefined if there aren't enough
// bits.
func (bs *Bitstring) Gray64(val int) uint64 {
	v := bs.Uint64(val)
	v ^= v >> 32
	v ^= v >> 16
	v ^= v >> 8
	v ^= v >> 4
	v ^= v >> 2
	v ^= v >> 1
	return v
}

// Grayn interprets the n bits at offset off as an n-bit gray-coded signed
// integer in big endian and returns its value. Behavior is undefined if there
// aren't enough bits. Panics if nbits is greater than 64.
func (bs *Bitstring) Grayn(off, n int) uint64 {
	v := bs.Uintn(off, n)
	v ^= v >> 32
	v ^= v >> 16
	v ^= v >> 8
	v ^= v >> 4
	v ^= v >> 2
	v ^= v >> 1
	return v
}
