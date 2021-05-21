package bitstring

// bitmask returns a mask where only the nth bit of a uint is set.
func bitmask(n uint64) uint64 { return 1 << n }

// wordoffset returns, for a given bit n of a bit string, the offset
// of the uint64 that contains bit n.
func wordoffset(n uint64) uint64 { return n / 64 }

// bitoffset returns, for a given bit n of a bit string, the offset of that bit
// with respect to the first bit of the uint64 that contains it.
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

// transferbits returns the uint that results from transferring some bits from
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
