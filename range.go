package bitstring

func minuint(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

// SwapRange swaps a range of bits between 2 bitstrings.
//
// The range [off, off+len) must exist on both bitstrings or SwapRange has
// undefined behavior.
func SwapRange(bs1, bs2 *Bitstring, off, len int) {
	bs1.mustExist(off + len - 1)
	bs2.mustExist(off + len - 1)

	// Swap bits in the first word.
	start, l := uint64(off), uint64(len)
	i := wordoffset(start)
	start = bitoffset(start)
	end := minuint(start+l, uintsize)
	swapBits(bs1, bs2, i, mask(start, end))
	i++

	// Swap whole words but the last one.
	remain := l - (end - start)
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
// It's like Equals but only compares the [off, off+length) range. EqualRange
// returns false if this range is not defined on both bitstrings.
func EqualRange(bs1, bs2 *Bitstring, off, len int) bool {
	if off+len-1 >= bs1.length || off+len-1 >= bs2.length {
		return false
	}

	// Compare bits in the first word.
	start, l := uint64(off), uint64(len)
	i := wordoffset(start)
	start = bitoffset(start)
	end := minuint(start+l, uintsize)
	m := mask(start, end)
	if bs1.data[i]&m != bs2.data[i]&m {
		return false
	}
	i++

	// Compare all words but the last one.
	remain := l - (end - start)
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
// The range [off, off+len) must exist or SetBitRange has undefined behavior.
func (bs *Bitstring) SetRange(off, len int) {
	bs.mustExist(off + len - 1)

	// Set bits in the first word.
	start, l := uint64(off), uint64(len)
	i := wordoffset(start)
	start = bitoffset(start)
	end := minuint(start+l, uintsize)
	bs.data[i] |= mask(start, end)
	i++

	// Set bits in all words but the last one.
	remain := l - (end - start)
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
// The range [off, off+length) must exist or ClearRange has undefined behavior.
func (bs *Bitstring) ClearRange(off, len int) {
	bs.mustExist(off + len - 1)

	// Clear bits in the first word.
	start, l := uint64(off), uint64(len)
	i := wordoffset(start)
	start = bitoffset(start)
	end := minuint(start+l, uintsize)
	bs.data[i] &= ^mask(start, end)
	i++

	// Clear bits in all words but the last one.
	remain := l - (end - start)
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
// The range [off, off+len) must exist or FlipRange has undefined behavior.
func (bs *Bitstring) FlipRange(off, len int) {
	bs.mustExist(off + len - 1)

	// Flip bits in the first word.
	start, l := uint64(off), uint64(len)
	i := wordoffset(start)
	start = bitoffset(start)
	end := minuint(start+l, uintsize)
	bs.data[i] ^= mask(start, end)
	i++

	// Flip bits in all words but the last one.
	remain := l - (end - start)
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
