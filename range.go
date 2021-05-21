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
func SwapRange(bs1, bs2 *Bitstring, off, len int) {
	bs1.mustExist(off + len - 1)
	bs2.mustExist(off + len - 1)

	// Swap bits in the first word.
	start64, len64 := uint64(off), uint64(len)
	i := wordoffset(start64)
	start64 = bitoffset(start64)
	end := minuint(start64+len64, uintsize)
	swapBits(bs1, bs2, i, mask(start64, end))
	i++

	// Swap whole words but the last one.
	remain := len64 - (end - start64)
	for remain > uintsize {
		bs1.data[i], bs2.data[i] = bs2.data[i], bs1.data[i]
		remain -= uintsize
		i++
	}

	// Swap bits in the last word.
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

// EqualRange compares a given range of bits between 2 bitstrings.
//
// It's like Equals but only compares the [start, start+length) range.
// EqualRange returns false if this range is not defined on both bitstrings.
func EqualRange(bs1, bs2 *Bitstring, start, length int) bool {
	if start+length-1 >= bs1.length || start+length-1 >= bs2.length {
		return false
	}

	// Compare bits in the first word.
	start64, len64 := uint64(start), uint64(length)
	i := wordoffset(start64)
	start64 = bitoffset(start64)
	end := minuint(start64+len64, uintsize)
	m := mask(start64, end)
	if bs1.data[i]&m != bs2.data[i]&m {
		return false
	}
	i++

	// Compare all words but the last one.
	remain := len64 - (end - start64)
	j := i + (remain / uintsize)
	if !u64cmp(bs1.data[i:j], bs2.data[i:j]) {
		return false
	}

	// Compare bits in the last word.
	if remain != 0 {
		m := lomask(remain)
		if bs1.data[j]&m != bs2.data[j]&m {
			return false
		}
	}

	return true
}

// SetRange sets a range of bits (sets all bits to 1).
//
// The range [start, start+length) must exist or SetBitRange has undefined
// behavior.
func (bs *Bitstring) SetRange(start, length int) {
	bs.mustExist(start + length - 1)

	// Sets bits in the first word.
	start64, len64 := uint64(start), uint64(length)
	i := wordoffset(start64)
	start64 = bitoffset(start64)
	end := minuint(start64+len64, uintsize)
	bs.data[i] |= mask(start64, end)
	i++

	// Set bits in all words but the last one.
	remain := len64 - (end - start64)
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

// ClearRange clears a range of bits (sets all bits to 0).
//
// The range [start, start+length) must exist or ClearRange has undefined
// behavior.
func (bs *Bitstring) ClearRange(start, length int) {
	bs.mustExist(start + length - 1)

	// Clears bits in the first word.
	start64, len64 := uint64(start), uint64(length)
	i := wordoffset(start64)
	start64 = bitoffset(start64)
	end := minuint(start64+len64, uintsize)
	bs.data[i] &= ^mask(start64, end)
	i++

	// Clears bits in all words but the last one.
	remain := len64 - (end - start64)
	for remain > uintsize {
		bs.data[i] = 0
		remain -= uintsize
		i++
	}

	// Clear bits in the last word.
	if remain != 0 {
		bs.data[i] &= himask(remain)
	}
}

// FlipRange flips a range of bits (flips the value of every bit).
//
// The range [start, start+length) must exist or FlipRange has undefined
// behavior.
func (bs *Bitstring) FlipRange(start, length int) {
	bs.mustExist(start + length - 1)

	// Flips bits in the first word.
	start64, len64 := uint64(start), uint64(length)
	i := wordoffset(start64)
	start64 = bitoffset(start64)
	end := minuint(start64+len64, uintsize)
	bs.data[i] ^= mask(start64, end)
	i++

	// Flips bits in all words but the last one.
	remain := len64 - (end - start64)
	for remain > uintsize {
		bs.data[i] ^= maxuint
		remain -= uintsize
		i++
	}

	// Flip bits in the last word.
	if remain != 0 {
		bs.data[i] ^= lomask(remain)
	}
}
