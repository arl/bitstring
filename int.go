package bitstring

/* get */

// Int8 returns the int8 value represented by the 8 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Int8(i int) int8 { return int8(bs.Uint8(i)) }

// Int16 returns the int16 value represented by the 16 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Int16(i int) int16 { return int16(bs.Uint16(i)) }

// Int32 returns the int32 value represented by the 32 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Int32(i int) int32 { return int32(bs.Uint32(i)) }

// Intn returns the n-bit signed integer value represented by the n bits
// starting at the i. It panics if there are not enough bits or if n is greater
// than the size of a machine word.
func (bs *Bitstring) Intn(nbits, i int) int32 { return int32(bs.Uintn(nbits, i)) }

// Int64 returns the int64 value represented by the 64 bits starting at the
// given bit. It panics if there are not enough bits.
func (bs *Bitstring) Int64(i int) int64 { return int64(bs.Uint64(i)) }

/* set */

// SetInt8 sets the 8 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetInt8(i int, x int8) { bs.SetUint8(i, uint8(x)) }

// SetInt16 sets the 16 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetInt16(i int, x int16) { bs.SetUint16(i, uint16(x)) }

// SetInt32 sets the 32 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetInt32(i int, x int32) { bs.SetUint32(i, uint32(x)) }

// SetInt64 sets the 64 bits starting at i with the value of x. It panics if
// there are not enough bits.
func (bs *Bitstring) SetInt64(i int, x int64) { bs.SetUint64(i, uint64(x)) }

// SetIntn sets the n bits starting at i with the first n bits of value x. It
// panics if there aren't enough bits in bs or if n is greater than 64.
func (bs *Bitstring) SetIntn(n, i int, x int64) { bs.SetUintn(n, i, uint64(x)) }
