package bitstring

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test8bits(t *testing.T) {
	var tests = []struct {
		str  string
		off  int
		want uint8
	}{
		// Value lies on a single uint64
		{str: "00000000000000000000000000000001", off: 0, want: 1},
		{str: "00000000000000000000000000000010", off: 0, want: 2},
		{str: "000000000000000001000010", off: 0, want: 1<<6 + 2},
		{str: "111111111111111101000010", off: 0, want: 1<<6 + 2},
		{str: "0000000000000000000000000000000111111111111111111111111111111111", off: 32, want: 1},
		{str: "00000000000000001000000111111111111111111111111111111111", off: 32, want: 1<<7 + 1},
		{str: "100000000", off: 1, want: 1 << 7},
		// Value lies across 2 different uint64
		{str: "1111111111111111111111101000010111111111111111111111111111111111111111111111111111111111111111", off: 63, want: 1<<6 + 2},
		{str: "0000000000000000000000011111110000000000000000000000000000000011111111111111111111111111111111", off: 63, want: math.MaxUint8 - 1},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.str)
			u8 := bs.Uint8(tt.off)
			assert.Equalf(t, tt.want, u8, "got %s want %s",
				sprintubits(uint64(u8), 8), sprintubits(uint64(tt.want), 8))

			i8 := bs.Int8(tt.off)
			assert.Equalf(t, int8(tt.want), i8, "got %s want %s",
				sprintubits(uint64(i8), 8), sprintubits(uint64(tt.want), 8))
		})
	}
}

func Test16bits(t *testing.T) {
	var tests = []struct {
		str  string
		off  int
		want uint16
	}{
		// Value lies on a single uint64
		{str: "00000000000000000000000000000001", off: 0, want: 1},
		{str: "00000000000000000000000000000010", off: 0, want: 2},
		{str: "00000000000000000100000000000010", off: 0, want: 1<<14 + 2},
		{str: "11111111111111110100000000000010", off: 0, want: 1<<14 + 2},
		{str: "0000000000000000000000000000000111111111111111111111111111111111", off: 32, want: 1},
		{str: "0000000000000000100000000000000111111111111111111111111111111111", off: 32, want: 1<<15 + 1},
		{str: "10000000000000000", off: 1, want: 1 << 15},
		// Value lies across 2 different uint64
		{str: "11111111111111111111111010000000000001011111111111111111111111111111111111111111111111111111111", off: 56, want: 1<<14 + 2},
		{str: "00000000000000000000000111111111111111000000000000000000000000011111111111111111111111111111111", off: 56, want: math.MaxUint16 - 1},
	}
	for _, tt := range tests {
		bs, _ := NewFromString(tt.str)
		u16 := bs.Uint16(tt.off)
		assert.Equalf(t, tt.want, u16, "got %s want %s",
			sprintubits(uint64(u16), 16), sprintubits(uint64(tt.want), 16))

		i16 := bs.Int16(tt.off)
		assert.Equalf(t, int16(tt.want), i16, "got %s want %s",
			sprintubits(uint64(i16), 16), sprintubits(uint64(tt.want), 16))
	}
}

func Test32bits(t *testing.T) {
	var tests = []struct {
		str  string
		off  int
		want uint32
	}{
		// Value lies on a single uint64
		{str: "00000000000000000000000000000001", off: 0, want: 1},
		{str: "00000000000000000000000000000010", off: 0, want: 2},
		{str: "01000000000000000000000000000010", off: 0, want: 1<<30 + 2},
		{str: "1111111111111111111111111111111101000000000000000000000000000010", off: 0, want: 1<<30 + 2},
		{str: "0000000000000000000000000000000111111111111111111111111111111111", off: 32, want: 1},
		{str: "1000000000000000000000000000000111111111111111111111111111111111", off: 32, want: 1<<31 + 1},
		{str: "100000000000000000000000000000000", off: 1, want: 1 << 31},
		// Value lies across 2 different uint64
		{str: "1111111111111111111101000000000000000000000000000010111111111111111111111111111111111111111111111111111111111111", off: 60, want: 1<<30 + 2},
		{str: "000011111111111111111111111111111110000000000000000000000000000011111111111111111111111111111111", off: 60, want: math.MaxUint32 - 1},
	}
	for _, tt := range tests {
		bs, _ := NewFromString(tt.str)
		u32 := bs.Uint32(tt.off)
		assert.Equalf(t, tt.want, u32, "got %s want %s",
			sprintubits(uint64(u32), 32), sprintubits(uint64(tt.want), 32))

		i32 := bs.Int32(tt.off)
		assert.Equalf(t, int32(tt.want), i32, "got %s want %s",
			sprintubits(uint64(i32), 32), sprintubits(uint64(tt.want), 32))
	}
}

func Test64bits(t *testing.T) {
	var tests = []struct {
		str  string
		off  int
		want uint64
	}{
		// Value lies on a single uint64
		{str: "0000000000000000000000000000000000000000000000000000000000000001", off: 0, want: 1},
		{str: "0000000000000000000000000000000000000000000000000000000000000010", off: 0, want: 2},
		{str: "0100000000000000000000000000000000000000000000000000000000000010", off: 0, want: 1<<62 + 2},
		{str: "11111111111111111111111111111111111111111111111111111111111111110100000000000000000000000000000000000000000000000000000000000010", off: 0, want: 1<<62 + 2},
		{str: "00000000000000000000000000000000000000000000000000000000000000011111111111111111111111111111111111111111111111111111111111111111", off: 64, want: 1},
		{str: "10000000000000000000000000000000000000000000000000000000000000011111111111111111111111111111111111111111111111111111111111111111", off: 64, want: 1<<63 + 1},
		{str: "10000000000000000000000000000000000000000000000000000000000000000", off: 1, want: 1 << 63},
		// Value lies across 2 different uint64
		{str: "1111111111111111111111111110100000000000000000000000000000000000000000000000000000000000010111111111111111111111111111111111111111111111111111111111111", off: 60, want: 1<<62 + 2},
		{str: "000011111111111111111111111111111111111111111111111111111111111111100000000000000000000000000000000000000000000000000000000000", off: 58, want: math.MaxUint64 - 1},
	}
	for _, tt := range tests {
		bs, _ := NewFromString(tt.str)
		u64 := bs.Uint64(tt.off)
		assert.Equalf(t, tt.want, u64, "got %s want %s",
			sprintubits(uint64(u64), 64), sprintubits(uint64(tt.want), 64))

		i64 := bs.Int64(tt.off)
		assert.Equalf(t, int64(tt.want), i64, "got %s want %s",
			sprintubits(uint64(i64), 64), sprintubits(uint64(tt.want), 64))
	}
}

func TestNbits(t *testing.T) {
	var tests = []struct {
		str    string
		n, off int
		want   uint64
	}{
		// Value lies on a single uint64
		{str: "10", n: 1, off: 0, want: 0},
		{str: "111", n: 1, off: 0, want: 1},
		{str: "101", n: 1, off: 1, want: 0},
		{str: "010", n: 1, off: 1, want: 1},
		{str: "100", n: 2, off: 0, want: 0},
		{str: "1101", n: 2, off: 1, want: 2},
		{str: "10100000000000000000000000000000", n: 3, off: 29, want: 5},
		{str: "10000000000000000000000000000000", n: 1, off: 31, want: 1},
		{str: "1111111111111111111111111111111111111111111111111111111111111111", n: 3, off: 31, want: 7},
		{str: "1111111111111111111111111111111111111111111111111111111111111111", n: 3, off: 30, want: 7},
		{str: "0000000000000000000000000000001010000000000000000000000000000000", n: 3, off: 31, want: 5},
		{str: "0000000000000000000000000000000101000000000000000000000000000000", n: 3, off: 30, want: 5},
		// Value lies across 2 different uint64
		{str: "000000000000000000000000000000101000000000000000000000000000000000000000000000000000000000000000", n: 3, off: 63, want: 5},
		{str: "000000000000000000000000000000010100000000000000000000000000000000000000000000000000000000000000", n: 3, off: 62, want: 5},
	}
	for _, tt := range tests {
		bs, _ := NewFromString(tt.str)
		un := bs.Uintn(tt.off, tt.n)
		assert.Equalf(t, tt.want, un, "got %s want %s",
			sprintubits(uint64(un), tt.n), sprintubits(uint64(tt.want), tt.n))

		in := bs.Intn(tt.off, tt.n)
		assert.Equalf(t, int64(tt.want), in, "got %s want %s",
			sprintubits(uint64(in), tt.n), sprintubits(uint64(tt.want), tt.n))
	}
}

func TestSet8bits(t *testing.T) {
	tests := []struct {
		str  string // starting bitstring
		val  uint8  // value to set
		off  int    // bit offset where to set it
		want string
	}{
		// Value lies on a single uint64
		{
			val: 2, off: 0,
			str:  "1111111111111111",
			want: "1111111100000010",
		},
		{
			val: 2, off: 8,
			str:  "1111111111111111",
			want: "0000001011111111",
		},
		{
			val: 2, off: 16,
			str:  "11111111111111111111111111111111",
			want: "11111111000000101111111111111111",
		},
		{
			val: 2, off: 24,
			str:  "11111111111111111111111111111111",
			want: "00000010111111111111111111111111",
		},
		{
			val: 2, off: 22,
			str:  "11111111111111111111111111111111",
			want: "11000000101111111111111111111111",
		},
		// Value lies across 2 different uint64
		{
			val: 2, off: 61,
			str:  "111111111111111111111111111111111111111111111111111111111111111111111",
			want: "000000101111111111111111111111111111111111111111111111111111111111111",
		},
		{
			val: 15, off: 63,
			str:  "11111111111111111111111111111111111111111111111111111111111111111111111",
			want: "00001111111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			val: math.MaxUint8, off: 59,
			str:  "0000000000000000000000000000000000000000000000000000000000000000000",
			want: "1111111100000000000000000000000000000000000000000000000000000000000",
		},
		{
			val: 0xaa, off: 63,
			str:  "0011101010101010101010101010101010101010101010101010101010101010101010101010101010",
			want: "0011101010110101010010101010101010101010101010101010101010101010101010101010101010",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.str)
			bs.SetUint8(tt.off, tt.val)
			assert.Equal(t, tt.want, bs.String())

			bs.SetInt8(tt.off, int8(tt.val))
			assert.Equal(t, tt.want, bs.String())
		})
	}
}

func TestSet16bits(t *testing.T) {
	tests := []struct {
		str  string // starting bitstring
		val  uint16 // value to set
		off  int    // bit offset where to set it
		want string
	}{
		// Value lies on a single uint64
		{
			val: 2, off: 0,
			str:  "1111111111111111",
			want: "0000000000000010",
		},
		{
			val: 2, off: 8,
			str:  "111111111111111111111111",
			want: "000000000000001011111111",
		},
		{
			val: 2, off: 16,
			str:  "11111111111111111111111111111111",
			want: "00000000000000101111111111111111",
		},
		{
			val: 2, off: 24,
			str:  "1111111111111111111111111111111111111111",
			want: "0000000000000010111111111111111111111111",
		},
		{
			val: 2, off: 22,
			str:  "1111111111111111111111111111111111111111",
			want: "1100000000000000101111111111111111111111",
		},
		// Value lies across 2 different uint64
		{
			val: 2, off: 61,
			str:  "11111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "00000000000000101111111111111111111111111111111111111111111111111111111111111",
		},
		{
			val: 15, off: 63,
			str:  "1111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "0000000000001111111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			val: math.MaxUint16, off: 59,
			str:  "000000000000000000000000000000000000000000000000000000000000000000000000000",
			want: "111111111111111100000000000000000000000000000000000000000000000000000000000",
		},
		{
			val: 0xaaaa, off: 63,
			str:  "001110101010101010101010101010101010101010101010101010101010101010101010101010101010101010",
			want: "001110101011010101010101010010101010101010101010101010101010101010101010101010101010101010",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.str)
			bs.SetUint16(tt.off, tt.val)
			assert.Equal(t, tt.want, bs.String())

			bs.SetInt16(tt.off, int16(tt.val))
			assert.Equal(t, tt.want, bs.String())
		})
	}
}

func TestSet32bits(t *testing.T) {
	tests := []struct {
		str  string // starting bitstring
		val  uint32 // value to set
		off  int    // bit offset where to set it
		want string
	}{
		// Value lies on a single uint64
		{
			val: 2, off: 0,
			str:  "11111111111111111111111111111111",
			want: "00000000000000000000000000000010",
		},
		{
			val: 2, off: 8,
			str:  "1111111111111111111111111111111111111111",
			want: "0000000000000000000000000000001011111111",
		},
		{
			val: 2, off: 16,
			str:  "111111111111111111111111111111111111111111111111",
			want: "000000000000000000000000000000101111111111111111",
		},
		{
			val: 2, off: 24,
			str:  "11111111111111111111111111111111111111111111111111111111",
			want: "00000000000000000000000000000010111111111111111111111111",
		},
		{
			val: 2, off: 22,
			str:  "11111111111111111111111111111111111111111111111111111111",
			want: "11000000000000000000000000000000101111111111111111111111",
		},
		// Value lies across 2 different uint64
		{
			val: 2, off: 61,
			str:  "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "000000000000000000000000000000101111111111111111111111111111111111111111111111111111111111111",
		},
		{
			val: 15, off: 63,
			str:  "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "00000000000000000000000000001111111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			val: math.MaxUint32, off: 59,
			str:  "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			want: "1111111111111111111111111111111100000000000000000000000000000000000000000000000000000000000",
		},
		{
			val: 0xaaaaaaaa, off: 63,
			str:  "0011101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010",
			want: "0011101010110101010101010101010101010101010010101010101010101010101010101010101010101010101010101010101010",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.str)
			bs.SetUint32(tt.off, tt.val)
			assert.Equal(t, tt.want, bs.String())

			bs.SetInt32(tt.off, int32(tt.val))
			assert.Equal(t, tt.want, bs.String())
		})
	}
}

func TestSet64bits(t *testing.T) {
	tests := []struct {
		str  string // starting bitstring
		val  uint64 // value to set
		off  int    // bit offset where to set it
		want string
	}{
		// Value lies on a single uint64
		{
			val: 2, off: 0,
			str:  "1111111111111111111111111111111111111111111111111111111111111111",
			want: "0000000000000000000000000000000000000000000000000000000000000010",
		},
		// Value lies across 2 different uint64
		{
			val: 2, off: 8,
			str:  "111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "000000000000000000000000000000000000000000000000000000000000001011111111",
		},
		{
			val: 2, off: 16,
			str:  "11111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "00000000000000000000000000000000000000000000000000000000000000101111111111111111",
		},
		{
			val: 2, off: 24,
			str:  "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "0000000000000000000000000000000000000000000000000000000000000010111111111111111111111111",
		},
		{
			val: 2, off: 22,
			str:  "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "1100000000000000000000000000000000000000000000000000000000000000101111111111111111111111",
		},
		{
			val: 2, off: 61,
			str:  "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "00000000000000000000000000000000000000000000000000000000000000101111111111111111111111111111111111111111111111111111111111111",
		},
		{
			val: 15, off: 63,
			str:  "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "0000000000000000000000000000000000000000000000000000000000001111111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			val: math.MaxUint64, off: 59,
			str:  "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			want: "111111111111111111111111111111111111111111111111111111111111111100000000000000000000000000000000000000000000000000000000000",
		},
		{
			val: 0xaaaaaaaaaaaaaaaa, off: 63,
			str:  "001110101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010",
			want: "001110101011010101010101010101010101010101010101010101010101010101010101010010101010101010101010101010101010101010101010101010101010101010",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.str)
			bs.SetUint64(tt.off, tt.val)
			assert.Equal(t, tt.want, bs.String())

			bs.SetInt64(tt.off, int64(tt.val))
			assert.Equal(t, tt.want, bs.String())
		})
	}
}

func TestSetNbits(t *testing.T) {
	tests := []struct {
		str    string // starting bitstring
		val    uint64 // value to set
		n, off int
		want   string
	}{
		// Value lies on a single uint64
		{
			n: 2, val: 2, off: 0,
			str:  "000",
			want: "010",
		},
		{
			n: 2, val: 2, off: 1,
			str:  "000",
			want: "100",
		},
		{
			n: 2, val: 2, off: 1,
			str:  "1111",
			want: "1101",
		},
		{
			n: 1, val: 1, off: 19,
			str:  "00000000000000000000000000000000",
			want: "00000000000010000000000000000000",
		},
		{
			n: 1, val: 3, off: 19,
			str:  "00000000000000000000000000000000",
			want: "00000000000010000000000000000000",
		},
		{
			n: 4, val: 8, off: 19,
			str:  "00000000011110000000000000000000",
			want: "00000000010000000000000000000000",
		},
		{
			n: 32, val: 0x80000002, off: 4,
			str:  "0101010101010101010101010101010101010101",
			want: "0101100000000000000000000000000000100101",
		},
		// Value lies across 2 different uint64
		{
			n: 2, val: 3, off: 63,
			str:  "0000000000000000000000000000000000000000000000000000000000000000000",
			want: "0011000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			n: 4, val: 9, off: 63,
			str:  "00000000000000000000000000000000000000000000000000000000000000000000",
			want: "01001000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			n: 64, val: 0x9cfbeb71ee3fcf5f, off: 35,
			str:  "000000000000000000001101000011010011000001010011010101010101000100101000111101010100000000000000000000000000000000000",
			want: "000000000000000000100111001111101111101011011100011110111000111111110011110101111100000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.str)
			bs.SetUintn(tt.off, tt.n, tt.val)
			assert.Equal(t, tt.want, bs.String())

			bs.SetIntn(tt.off, tt.n, int64(tt.val))
			assert.Equal(t, tt.want, bs.String())
		})
	}
}
