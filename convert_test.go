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
			got := bs.Uintn(tt.i, tt.nbits)
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
