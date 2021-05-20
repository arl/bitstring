package bitstring

import (
	"math"
	"testing"
)

//
// conversion to unsigned integers
//

func TestUintn(t *testing.T) {
	tests := []struct {
		input    string
		nbits, i int
		want     uint64
	}{
		// LSB and MSB are both on the same word
		{input: "10", nbits: 1, i: 0, want: 0},
		{input: "111", nbits: 1, i: 0, want: 1},
		{input: "101", nbits: 1, i: 1, want: 0},
		{input: "010", nbits: 1, i: 1, want: 1},
		{input: "100", nbits: 2, i: 0, want: 0},
		{input: "1101", nbits: 2, i: 1, want: 2},
		{input: "10100000000000000000000000000000", nbits: 3, i: 29, want: 5},
		{input: "10000000000000000000000000000000", nbits: 1, i: 31, want: 1},

		// // LSB and MSB are on 2 separate words
		{input: "1111111111111111111111111111111111111111111111111111111111111111", nbits: 3, i: 31, want: 7},
		{input: "1111111111111111111111111111111111111111111111111111111111111111", nbits: 3, i: 30, want: 7},
		{input: "0000000000000000000000000000001010000000000000000000000000000000", nbits: 3, i: 31, want: 5},
		{input: "0000000000000000000000000000000101000000000000000000000000000000", nbits: 3, i: 30, want: 5},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.input)
			got := bs.Uintn(tt.nbits, tt.i)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Uintn(%d, %d) got %s, want %s", tt.input, tt.nbits, tt.i,
					sprintubits(got, tt.nbits), sprintubits(tt.want, tt.nbits))
			}
		})
	}
}

func TestUint64(t *testing.T) {
	tests := []struct {
		input string
		i     int
		want  uint64
	}{
		// LSB and MSB are both on the same word
		{
			input: "0000000000000000000000000000000000000000000000000000000000000001",
			i:     0, want: 1,
		},
		{
			input: "0000000000000000000000000000000000000000000000000000000000000010",
			i:     0, want: 2,
		},
		{
			input: "0100000000000000000000000000000000000000000000000000000000000010",
			i:     0, want: 1<<62 + 2,
		},
		{
			input: "11111111111111111111111111111111111111111111111111111111111111110100000000000000000000000000000000000000000000000000000000000010",
			i:     0, want: 1<<62 + 2,
		},
		{
			input: "00000000000000000000000000000000000000000000000000000000000000011111111111111111111111111111111111111111111111111111111111111111",
			i:     64, want: 1,
		},
		{
			input: "10000000000000000000000000000000000000000000000000000000000000011111111111111111111111111111111111111111111111111111111111111111",
			i:     64, want: 1<<63 + 1,
		},

		// LSB and MSB are on 2 separate words
		{
			input: "10000000000000000000000000000000000000000000000000000000000000000",
			i:     1, want: 1 << 63,
		},
		{
			input: "1111111111111111111111111110100000000000000000000000000000000000000000000000000000000000010111111111111111111111111111111111111111111111111111111111111",
			i:     60, want: 1<<62 + 2,
		},
		{
			input: "000011111111111111111111111111111111111111111111111111111111111111100000000000000000000000000000000000000000000000000000000000",
			i:     58, want: math.MaxUint64 - 1,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.input)
			got := bs.Uint64(tt.i)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Uint64(%d) got %s, want %s", tt.input, tt.i,
					sprintubits(got, 64), sprintubits(tt.want, 64))
			}
		})
	}
}

func TestUint32(t *testing.T) {
	tests := []struct {
		input string
		i     int
		want  uint32
	}{
		// LSB and MSB are both on the same word
		{
			input: "00000000000000000000000000000001",
			i:     0, want: 1,
		},
		{
			input: "00000000000000000000000000000010",
			i:     0, want: 2,
		},
		{
			input: "01000000000000000000000000000010",
			i:     0, want: 1<<30 + 2,
		},
		{
			input: "1111111111111111111111111111111101000000000000000000000000000010",
			i:     0, want: 1<<30 + 2,
		},
		{
			input: "0000000000000000000000000000000111111111111111111111111111111111",
			i:     32, want: 1,
		},
		{
			input: "1000000000000000000000000000000111111111111111111111111111111111",
			i:     32, want: 1<<31 + 1,
		},

		// LSB and MSB are on 2 separate words
		{
			input: "100000000000000000000000000000000",
			i:     1, want: 1 << 31,
		},
		{
			input: "1111111111111111111101000000000000000000000000000010111111111111",
			i:     12, want: 1<<30 + 2,
		},
		{
			input: "0000111111111111111111111111111111100000000000000000000000000000",
			i:     28, want: math.MaxUint32 - 1,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.input)
			got := bs.Uint32(tt.i)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Uint32(%d) got %s, want %s", tt.input, tt.i,
					sprintubits(uint64(got), 32), sprintubits(uint64(tt.want), 32))
			}
		})
	}
}

func TestUint16(t *testing.T) {
	tests := []struct {
		input string
		i     int
		want  uint16
	}{
		// LSB and MSB are both on the same word
		{
			input: "00000000000000000000000000000001",
			i:     0, want: 1,
		},
		{
			input: "00000000000000000000000000000010",
			i:     0, want: 2,
		},
		{
			input: "00000000000000000100000000000010",
			i:     0, want: 1<<14 + 2,
		},
		{
			input: "11111111111111110100000000000010",
			i:     0, want: 1<<14 + 2,
		},
		{
			input: "0000000000000000000000000000000111111111111111111111111111111111",
			i:     32, want: 1,
		},
		{
			input: "0000000000000000100000000000000111111111111111111111111111111111",
			i:     32, want: 1<<15 + 1,
		},
		{
			input: "10000000000000000",
			i:     1, want: 1 << 15,
		},

		// LSB and MSB are on 2 separate words
		{
			input: "111111111111111111111110100000000000010111111111111111111111111",
			i:     24, want: 1<<14 + 2,
		},
		{
			input: "000000000000000000000001111111111111110000000000000000000000000",
			i:     24, want: math.MaxUint16 - 1,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.input)
			got := bs.Uint16(tt.i)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Uint16(%d) got %s, want %s", tt.input, tt.i,
					sprintubits(uint64(got), 16), sprintubits(uint64(tt.want), 16))
			}
		})
	}
}

func TestUint8(t *testing.T) {
	tests := []struct {
		input string
		i     int
		want  uint8
	}{
		// LSB and MSB are both on the same word
		{
			input: "00000000000000000000000000000001",
			i:     0, want: 1,
		},
		{
			input: "00000000000000000000000000000010",
			i:     0, want: 2,
		},
		{
			input: "000000000000000001000010",
			i:     0, want: 1<<6 + 2,
		},
		{
			input: "111111111111111101000010",
			i:     0, want: 1<<6 + 2,
		},
		{
			input: "0000000000000000000000000000000111111111111111111111111111111111",
			i:     32, want: 1,
		},
		{
			input: "00000000000000001000000111111111111111111111111111111111",
			i:     32, want: 1<<7 + 1,
		},
		{
			input: "100000000",
			i:     1, want: 1 << 7,
		},

		// LSB and MSB are on separate words
		{
			input: "11111111111111111111111010000101111111111111111111111111111111",
			i:     31, want: 1<<6 + 2,
		},
		{
			input: "00000000000000000000000111111100000000000000000000000000000000",
			i:     31, want: math.MaxUint8 - 1,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.input)
			got := bs.Uint8(tt.i)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Uint8(%d) got %s, want %s", tt.input, tt.i,
					sprintubits(uint64(got), 8), sprintubits(uint64(tt.want), 8))
			}
		})
	}
}

//
// conversion to signed integers
//
func TestInt32(t *testing.T) {
	tests := []struct {
		input string
		i     int
		want  int32
	}{
		// LSB and MSB are both on the same word
		{
			i: 0, want: -1,
			input: "11111111111111111111111111111111",
		},
		{
			i: 0, want: math.MaxInt32,
			input: "01111111111111111111111111111111",
		},
		{
			i: 0, want: math.MinInt32,
			input: "10000000000000000000000000000000",
		},
		// LSB and MSB are on 2 separate words
		{
			i: 31, want: -1,
			input: "111111111111111111111111111111110000000000000000000000000000000",
		},
		{
			i: 31, want: math.MaxInt32,
			input: "011111111111111111111111111111110000000000000000000000000000000",
		},
		{
			i: 31, want: math.MinInt32,
			input: "100000000000000000000000000000001111111111111111111111111111111",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.input)
			got := bs.Int32(tt.i)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Int32(%d) got %s, want %s", tt.input, tt.i,
					sprintsbits(int64(got), 32), sprintsbits(int64(tt.want), 32))
			}
		})
	}
}

func TestInt16(t *testing.T) {
	tests := []struct {
		input string
		i     int
		want  int16
	}{
		// LSB and MSB are both on the same word
		{
			i: 0, want: -1,
			input: "1111111111111111",
		},
		{
			i: 0, want: math.MaxInt16,
			input: "0111111111111111",
		},
		{
			i: 0, want: math.MinInt16,
			input: "1000000000000000",
		},
		// LSB and MSB are on 2 separate words
		{
			i: 31, want: -1,
			input: "11111111111111110000000000000000000000000000000",
		},
		{
			i: 31, want: math.MaxInt16,
			input: "01111111111111110000000000000000000000000000000",
		},
		{
			i: 31, want: math.MinInt16,
			input: "10000000000000001111111111111111111111111111111",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.input)
			got := bs.Int16(tt.i)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Int16(%d) got %s, want %s", tt.input, tt.i,
					sprintsbits(int64(got), 16), sprintsbits(int64(tt.want), 16))
			}
		})
	}
}

func TestInt8(t *testing.T) {
	tests := []struct {
		input string
		i     int
		want  int8
	}{
		// LSB and MSB are both on the same word
		{
			i: 0, want: -1,
			input: "11111111",
		},
		{
			i: 0, want: math.MaxInt8,
			input: "01111111",
		},
		{
			i: 0, want: math.MinInt8,
			input: "10000000",
		},
		// LSB and MSB are on 2 separate words
		{
			i: 31, want: -1,
			input: "111111110000000000000000000000000000000",
		},
		{
			i: 31, want: math.MaxInt8,
			input: "011111110000000000000000000000000000000",
		},
		{
			input: "100000001111111111111111111111111111111",
			i:     31, want: math.MinInt8,
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.input)
			got := bs.Int8(tt.i)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Int8(%d) got %s, want %s", tt.input, tt.i,
					sprintsbits(int64(got), 8), sprintsbits(int64(tt.want), 8))
			}
		})
	}
}

//
// conversion from unsigned integers
//

func TestSetUintn(t *testing.T) {
	tests := []struct {
		bs       string // starting bitstring
		x        uint64 // value to set
		nbits, i int
		want     string
	}{
		// LSB and MSB are both on the same word
		{
			nbits: 2, x: 2, i: 0,
			bs:   "000",
			want: "010",
		},
		{
			nbits: 2, x: 2, i: 1,
			bs:   "000",
			want: "100",
		},
		{
			nbits: 2, x: 2, i: 1,
			bs:   "1111",
			want: "1101",
		},
		{
			nbits: 1, x: 1, i: 19,
			bs:   "00000000000000000000000000000000",
			want: "00000000000010000000000000000000",
		},
		{
			nbits: 1, x: 3, i: 19,
			bs:   "00000000000000000000000000000000",
			want: "00000000000010000000000000000000",
		},
		{
			nbits: 4, x: 8, i: 19,
			bs:   "00000000011110000000000000000000",
			want: "00000000010000000000000000000000",
		},
		{
			nbits: 32, x: 0x80000002, i: 4,
			bs:   "0101010101010101010101010101010101010101",
			want: "0101100000000000000000000000000000100101",
		},
		// LSB and MSB are on separate words
		{
			nbits: 2, x: 3, i: 63,
			bs:   "0000000000000000000000000000000000000000000000000000000000000000000",
			want: "0011000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			nbits: 4, x: 9, i: 63,
			bs:   "00000000000000000000000000000000000000000000000000000000000000000000",
			want: "01001000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			nbits: 64, x: 0x9cfbeb71ee3fcf5f, i: 35,
			bs:   "000000000000000000001101000011010011000001010011010101010101000100101000111101010100000000000000000000000000000000000",
			want: "000000000000000000100111001111101111101011011100011110111000111111110011110101111100000000000000000000000000000000000",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.bs)
			bs.SetUintn(tt.nbits, tt.i, tt.x)
			want, _ := NewFromString(tt.want)
			if !want.Equals(bs) {
				t.Errorf("Bitstring(%s).SetUintn(%d, %d, %d) got %s, want %s",
					tt.bs, tt.nbits, tt.i, tt.x, bs, want)
			}
		})
	}
}

func TestSetUint8(t *testing.T) {
	tests := []struct {
		bs   string // starting bitstring
		x    uint8  // value to set
		i    int    // bit index where to set it
		want string
	}{
		// LSB and MSB are both on the same word
		{
			x: 2, i: 0,
			bs:   "1111111111111111",
			want: "1111111100000010",
		},
		{
			x: 2, i: 8,
			bs:   "1111111111111111",
			want: "0000001011111111",
		},
		{
			x: 2, i: 16,
			bs:   "11111111111111111111111111111111",
			want: "11111111000000101111111111111111",
		},
		{
			x: 2, i: 24,
			bs:   "11111111111111111111111111111111",
			want: "00000010111111111111111111111111",
		},
		{
			x: 2, i: 22,
			bs:   "11111111111111111111111111111111",
			want: "11000000101111111111111111111111",
		},
		// LSB and MSB are on separate words
		{
			x: 2, i: 61,
			bs:   "111111111111111111111111111111111111111111111111111111111111111111111",
			want: "000000101111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x: 15, i: 63,
			bs:   "11111111111111111111111111111111111111111111111111111111111111111111111",
			want: "00001111111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x: math.MaxUint8, i: 59,
			bs:   "0000000000000000000000000000000000000000000000000000000000000000000",
			want: "1111111100000000000000000000000000000000000000000000000000000000000",
		},
		{
			x: 0xaa, i: 63,
			bs:   "0011101010101010101010101010101010101010101010101010101010101010101010101010101010",
			want: "0011101010110101010010101010101010101010101010101010101010101010101010101010101010",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.bs)
			bs.SetUint8(tt.i, tt.x)
			want, _ := NewFromString(tt.want)
			if !want.Equals(bs) {
				t.Errorf("Bitstring(%s).SetUint8(%d, %d) got %s, want %s",
					tt.bs, tt.i, tt.x, bs, want)
			}
		})
	}
}

func TestSetUint16(t *testing.T) {
	tests := []struct {
		bs   string // starting bitstring
		x    uint16 // value to set
		i    int    // bit index where to set it
		want string
	}{
		// LSB and MSB are both on the same word
		{
			x: 2, i: 0,
			bs:   "1111111111111111",
			want: "0000000000000010",
		},
		{
			x: 2, i: 8,
			bs:   "111111111111111111111111",
			want: "000000000000001011111111",
		},
		{
			x: 2, i: 16,
			bs:   "11111111111111111111111111111111",
			want: "00000000000000101111111111111111",
		},
		{
			x: 2, i: 24,
			bs:   "1111111111111111111111111111111111111111",
			want: "0000000000000010111111111111111111111111",
		},
		{
			x: 2, i: 22,
			bs:   "1111111111111111111111111111111111111111",
			want: "1100000000000000101111111111111111111111",
		},
		// LSB and MSB are on separate words
		{
			x: 2, i: 61,
			bs:   "11111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "00000000000000101111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x: 15, i: 63,
			bs:   "1111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "0000000000001111111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x: math.MaxUint16, i: 59,
			bs:   "000000000000000000000000000000000000000000000000000000000000000000000000000",
			want: "111111111111111100000000000000000000000000000000000000000000000000000000000",
		},
		{
			x: 0xaaaa, i: 63,
			bs:   "001110101010101010101010101010101010101010101010101010101010101010101010101010101010101010",
			want: "001110101011010101010101010010101010101010101010101010101010101010101010101010101010101010",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.bs)
			bs.SetUint16(tt.i, tt.x)
			want, _ := NewFromString(tt.want)
			if !want.Equals(bs) {
				t.Errorf("Bitstring(%s).SetUint16(%d, %d) got %s, want %s",
					tt.bs, tt.i, tt.x, bs, want)
			}
		})
	}
}

func TestSetUint32(t *testing.T) {
	tests := []struct {
		bs   string // starting bitstring
		x    uint32 // value to set
		i    int    // bit index where to set it
		want string
	}{
		// LSB and MSB are both on the same word
		{
			x: 2, i: 0,
			bs:   "11111111111111111111111111111111",
			want: "00000000000000000000000000000010",
		},
		{
			x: 2, i: 8,
			bs:   "1111111111111111111111111111111111111111",
			want: "0000000000000000000000000000001011111111",
		},
		{
			x: 2, i: 16,
			bs:   "111111111111111111111111111111111111111111111111",
			want: "000000000000000000000000000000101111111111111111",
		},
		{
			x: 2, i: 24,
			bs:   "11111111111111111111111111111111111111111111111111111111",
			want: "00000000000000000000000000000010111111111111111111111111",
		},
		{
			x: 2, i: 22,
			bs:   "11111111111111111111111111111111111111111111111111111111",
			want: "11000000000000000000000000000000101111111111111111111111",
		},
		// LSB and MSB are on separate words
		{
			x: 2, i: 61,
			bs:   "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "000000000000000000000000000000101111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x: 15, i: 63,
			bs:   "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "00000000000000000000000000001111111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x: math.MaxUint32, i: 59,
			bs:   "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			want: "1111111111111111111111111111111100000000000000000000000000000000000000000000000000000000000",
		},
		{
			x: 0xaaaaaaaa, i: 63,
			bs:   "0011101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010",
			want: "0011101010110101010101010101010101010101010010101010101010101010101010101010101010101010101010101010101010",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.bs)
			bs.SetUint32(tt.i, tt.x)
			want, _ := NewFromString(tt.want)
			if !want.Equals(bs) {
				t.Errorf("Bitstring(%s).SetUint32(%d, %d) got %s, want %s",
					tt.bs, tt.i, tt.x, bs, want)
			}
		})
	}
}

func TestSetUint64(t *testing.T) {
	tests := []struct {
		bs   string // starting bitstring
		x    uint64 // value to set
		i    int    // bit index where to set it
		want string
	}{
		// LSB and MSB are both on the same word
		{
			x: 2, i: 0,
			bs:   "1111111111111111111111111111111111111111111111111111111111111111",
			want: "0000000000000000000000000000000000000000000000000000000000000010",
		},
		// LSB and MSB are on separate words
		{
			x: 2, i: 8,
			bs:   "111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "000000000000000000000000000000000000000000000000000000000000001011111111",
		},
		{
			x: 2, i: 16,
			bs:   "11111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "00000000000000000000000000000000000000000000000000000000000000101111111111111111",
		},
		{
			x: 2, i: 24,
			bs:   "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "0000000000000000000000000000000000000000000000000000000000000010111111111111111111111111",
		},
		{
			x: 2, i: 22,
			bs:   "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "1100000000000000000000000000000000000000000000000000000000000000101111111111111111111111",
		},
		{
			x: 2, i: 61,
			bs:   "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "00000000000000000000000000000000000000000000000000000000000000101111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x: 15, i: 63,
			bs:   "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			want: "0000000000000000000000000000000000000000000000000000000000001111111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x: math.MaxUint64, i: 59,
			bs:   "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			want: "111111111111111111111111111111111111111111111111111111111111111100000000000000000000000000000000000000000000000000000000000",
		},
		{
			x: 0xaaaaaaaaaaaaaaaa, i: 63,
			bs:   "001110101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010101010",
			want: "001110101011010101010101010101010101010101010101010101010101010101010101010010101010101010101010101010101010101010101010101010101010101010",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.bs)
			bs.SetUint64(tt.i, tt.x)
			want, _ := NewFromString(tt.want)
			if !want.Equals(bs) {
				t.Errorf("Bitstring(%s).SetUint64(%d, %d) got %s, want %s",
					tt.bs, tt.i, tt.x, bs, want)
			}
		})
	}
}

func TestSetInt(t *testing.T) {
	t.Run("SetInt8", func(t *testing.T) {
		bs, _ := NewFromString("00000000")
		want := int8(math.MinInt8 / 3 * 2)
		bs.SetInt8(0, want)
		got := bs.Int8(0)
		if got != want {
			t.Errorf("Bitstring().SetInt8(%d, %d) got %d, want %d",
				0, want, got, want)
		}
	})
	t.Run("SetInt16", func(t *testing.T) {
		bs, _ := NewFromString("0000000000000000")
		want := int16(math.MinInt16 / 3 * 2)
		bs.SetInt16(0, want)
		got := bs.Int16(0)
		if got != want {
			t.Errorf("Bitstring().SetInt16(%d, %d) got %d, want %d",
				0, want, got, want)
		}
	})
	t.Run("SetInt32", func(t *testing.T) {
		bs, _ := NewFromString("00000000000000000000000000000000")
		want := int32(math.MinInt32 / 3 * 2)
		bs.SetInt32(0, want)
		got := bs.Int32(0)
		if got != want {
			t.Errorf("Bitstring().SetInt32(%d, %d) got %d, want %d",
				0, want, got, want)
		}
	})
	t.Run("SetInt64", func(t *testing.T) {
		bs, _ := NewFromString("0000000000000000000000000000000000000000000000000000000000000000")
		want := int64(math.MinInt64 / 3 * 2)
		bs.SetInt64(0, want)
		got := bs.Int64(0)
		if got != want {
			t.Errorf("Bitstring().SetInt64(%d, %d) got %d, want %d",
				0, want, got, want)
		}
	})
}
