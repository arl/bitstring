package bitstring

func minuint(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

// SwapRange swaps a range of bits between 2 bitstrings.
//
// The range [start, start+length) must exist on both bitstrings or SwapRange
// has undefined behavior.
func SwapRange(bs1, bs2 *Bitstring, start, length int) {
	bs1.mustExist(start + length - 1)
	bs2.mustExist(start + length - 1)

	// Swap the required bits of the first word.
	start64, len64 := uint64(start), uint64(length)
	i := wordoffset(start64)
	start64 = bitoffset(start64)
	end := minuint(start64+len64, uintsize)
	remain := len64 - (end - start64)
	swapBits(bs1, bs2, i, mask(start64, end))
	i++

	// Swap whole words but the last one.
	for remain > uintsize {
		bs1.data[i], bs2.data[i] = bs2.data[i], bs1.data[i]
		remain -= uintsize
		i++
	}

	// Swap the remaining bits of the last word.
	if remain != 0 {
		swapBits(bs1, bs2, i, lomask(remain))
	}
}

// swapBits swaps range of bits from one word to another. w is the index of the
// word containing the bits to swap, and m is a mask indicating the bits to
// swap.
func swapBits(x, y *Bitstring, w, mask uint64) {
	keep := ^mask
	xkeep, ykeep := x.data[w]&keep, y.data[w]&keep
	xswap, yswap := x.data[w]&mask, y.data[w]&mask
	x.data[w] = xkeep | yswap
	y.data[w] = ykeep | xswap
}

// SetRange sets a range of bits (sets all bits to 1).
//
// The range [start, start+length) must exist or SetBitRange has undefined
// behavior.
func (bs *Bitstring) SetRange(start, length int) {
	bs.mustExist(start + length - 1)

	// Swap the required bits of the first word.
	start64, len64 := uint64(start), uint64(length)
	i := wordoffset(start64)
	start64 = bitoffset(start64)
	end := minuint(start64+len64, uintsize)
	remain := len64 - (end - start64)
	bs.data[i] |= mask(start64, end)
	i++

	// Set all bits in remaining words but the last one.
	for remain > uintsize {
		bs.data[i] = maxuint
		remain -= uintsize
		i++
	}

	// Set bits in the last word.
	if remain != 0 {
		bs.data[i] |= lomask(remain)
	}
}
