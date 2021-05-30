package bitstring

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	bs := New(100)

	assert.Equal(t, 100, bs.Len())
	for i := 0; i < bs.Len(); i++ {
		assert.False(t, bs.Bit(i))
	}
}

func TestRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(99))

	assert.Equal(t, Random(100, rng).Len(), 100)
}

func TestSetBit(t *testing.T) {
	bs := New(69)

	bs.SetBit(1)
	bs.SetBit(4)

	assert.False(t, bs.Bit(0))
	assert.True(t, bs.Bit(1))
	assert.False(t, bs.Bit(2))
	assert.False(t, bs.Bit(3))
	assert.True(t, bs.Bit(4))

	bs.ClearBit(4)
	assert.False(t, bs.Bit(4))
}

func TestFlipBit(t *testing.T) {
	bs := New(69)

	bs.FlipBit(2)
	assert.True(t, bs.Bit(2))

	bs.FlipBit(2)
	assert.False(t, bs.Bit(2))

	bs.FlipBit(67)
	assert.True(t, bs.Bit(67))

	bs.FlipBit(67)
	assert.False(t, bs.Bit(67))
}

func TestString(t *testing.T) {
	tests := []struct {
		name  string
		str   string
		valid bool
	}{
		// valid cases
		{"only 0s", "0000", true},
		{"only 1s", "1111111", true},
		{"mixed 0s and 1s", "1000111011", true},
		{"leading 0", "00001000111011", true},
		{"empty string", "", true},
		{"valid long", "111010101110101100010100101000101000011010100010111100011101001010101110101011010101101011100000000110110101011110101010101", true},

		// error cases
		{"with spaces", "11 ", false},
		{"invalid ascii chars", "0101012", false},
		{"non ascii chars", "10日本", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFromString(tt.str)
			if !tt.valid {
				assert.Nil(t, got)
				assert.Error(t, err)
				return
			}
			assert.NotNil(t, got)
			assert.NoError(t, err)

			// Convert back to string
			assert.Equal(t, tt.str, got.String())
		})
	}
}

func TestReverse(t *testing.T) {
	tests := []string{
		"0000000000000000000000000000000000000000000000000000000000000001",
		"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
		"1100010010111100010010111100010010110010101000101101001011111110011010010000100101101011100110100100001001011010100001101111010011000100101111000100101111000100101100101010001011010010111111100110100100001001011010111001101001000010010110101000011011110100",
		"00000001",
		"11110000000000000000000000000000000000000000000000000000000000000111",
		"10000000000000000000000000000000000000000000000000000000000000011111",
		fmt.Sprintf("%s000000000000000000", strings.Repeat("1", 1029)),
		"000000000000000000111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
	}
	for _, str := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(str)
			ones := bs.OnesCount()

			bs.Reverse()

			assert.Equal(t, string(reverseBytes([]byte(str))), bs.String())
			assert.Equal(t, ones, bs.OnesCount())
		})
	}
}

func TestOnesCount(t *testing.T) {
	bs := New(65)
	assert.Zero(t, bs.OnesCount())

	bs.SetBit(0)
	bs.SetBit(31)
	bs.SetBit(32)
	bs.SetBit(33)
	bs.SetBit(63)
	bs.SetBit(64)
	setBits := bs.OnesCount()
	assert.EqualValues(t, 6, setBits)
}

func TestZeroesCount(t *testing.T) {
	bs := New(12)
	assert.EqualValues(t, 12, bs.ZeroesCount())

	bs.SetBit(0)
	bs.SetBit(5)
	bs.SetBit(6)
	bs.SetBit(9)
	bs.SetBit(10)
	setBits := bs.ZeroesCount()
	assert.EqualValues(t, 7, setBits)
}

func TestClone(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	bs := Random(2000, rng)
	bs.SetBit(3)
	bs.SetBit(7)
	bs.SetBit(8)

	cpy := bs.Clone()
	equalbits(t, cpy, bs)
}

func TestCopy(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	bs := Random(2000, rng)
	bs.SetBit(3)
	bs.SetBit(7)
	bs.SetBit(1898)

	// zero-value destination
	var cpy Bitstring
	Copy(&cpy, bs)
	equalbits(t, bs, &cpy)

	// smaller destination
	cpy2 := New(1)
	Copy(cpy2, bs)
	equalbits(t, bs, cpy2)

	// same size destination
	cpy3 := New(2000)
	Copy(cpy3, bs)
	equalbits(t, bs, cpy3)

	// larger destination
	cpy4 := New(3000)
	Copy(cpy4, bs)
	equalbits(t, bs, cpy4)
}

func TestEquals(t *testing.T) {
	org := New(10)
	org.SetBit(2)
	org.SetBit(5)
	org.SetBit(8)

	assert.True(t, org.Equals(org))
	assert.False(t, org.Equals(nil))
	assert.False(t, org.Equals(&Bitstring{}))

	clone := org.Clone()
	assert.Truef(t, clone.Equals(org), "different bitstrings, clone=%s, org=%s", clone, org)

	clone.FlipBit(0)
	assert.Falsef(t, clone.Equals(org), "same bitstrings, clone=%s, org=%s", clone, org)

	// Bitstrings of different lengths but with the same bits set should not be equal.
	bs2 := New(9)
	bs2.SetBit(2)
	bs2.SetBit(5)
	bs2.SetBit(8)
	assert.False(t, bs2.Equals(org))
}

func trimLeadingZeroes(s string) string {
	last := len(s) - 1
	return strings.TrimLeft(s[:last], "0") + string(s[last])
}

func TestBig(t *testing.T) {
	nums := []string{
		"00101",
		"00101110",
		"001011100010111000101110001011100010111000101110001011100010111",
		"0010111000101110001011100010111000101110001011100010111000101110",
		"00101110001011100010111000101110001011100010111000101110001011101",
		"001011100010111000101110001011100010111000101110001011100010111001010101000111101010101000111101",
		"0010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111",
		"00101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110",
		"001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011101",
		"00101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111",
		"001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110",
		"0010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011100010111000101110001011101",
	}
	for _, num := range nums {
		t.Run(fmt.Sprintf("%d-bit", len(num)), func(t *testing.T) {
			// Create a bitstring, converts it to big.Int and format it as a
			// binary number, compare the numbers with the original.
			bs, _ := NewFromString(num)
			bi := bs.BigInt()
			sbig := fmt.Sprintf("%b", bi)
			assert.Equal(t, trimLeadingZeroes(num), sbig)

			// Create a bitstring from the binary string representation, create
			// a Bitstring from it and compare the significant bits.
			bs2 := NewFromBig(bi)

			// Remove leading zeroes since a big int doesn't have leading zeroes.
			got := bs.String()[bs.LeadingZeroes():]
			assert.Equal(t, got, bs2.String())
		})
	}
}

func TestLeadingTrailingZeroesOnes(t *testing.T) {
	tests := []struct {
		name                          string
		s                             string
		leadingZeroes, trailingZeroes int
		leadingOnes, trailingOnes     int
	}{
		{
			name:          "less than 64 bits",
			s:             "0",
			leadingZeroes: 1, trailingZeroes: 1,
			leadingOnes: 0, trailingOnes: 0,
		},
		{
			name:          "less than 64 bits",
			s:             "1",
			leadingZeroes: 0, trailingZeroes: 0,
			leadingOnes: 1, trailingOnes: 1,
		},
		{
			name:          "less than 64 bits",
			s:             "10",
			leadingZeroes: 0, trailingZeroes: 1,
			leadingOnes: 1, trailingOnes: 0,
		},
		{
			name:          "less than 64 bits",
			s:             "01",
			leadingZeroes: 1, trailingZeroes: 0,
			leadingOnes: 0, trailingOnes: 1,
		},
		{
			name:          "less than 64 bits",
			s:             "00",
			leadingZeroes: 2, trailingZeroes: 2,
			leadingOnes: 0, trailingOnes: 0,
		},
		{
			name:          "less than 64 bits",
			s:             "11",
			leadingZeroes: 0, trailingZeroes: 0,
			leadingOnes: 2, trailingOnes: 2,
		},
		{
			name:          "exactly 64 bits",
			s:             "1111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 0, trailingZeroes: 0,
			leadingOnes: 64, trailingOnes: 64,
		},
		{
			name:          "exactly 64 bits",
			s:             "0111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 1, trailingZeroes: 0,
			leadingOnes: 0, trailingOnes: 63,
		},
		{
			name:          "exactly 64 bits",
			s:             "1111111111111111111111111111111111111111111111111111111111111110",
			leadingZeroes: 0, trailingZeroes: 1,
			leadingOnes: 63, trailingOnes: 0,
		},
		{
			name:          "exactly 64 bits",
			s:             "1111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 0, trailingZeroes: 0,
			leadingOnes: 64, trailingOnes: 64,
		},
		{
			name:          "exactly 64 bits",
			s:             "0000000000000000000000000000000000000000000000000000000000000000",
			leadingZeroes: 64, trailingZeroes: 64,
			leadingOnes: 0, trailingOnes: 0,
		},
		{
			name:          "128 bits",
			s:             "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 0, trailingZeroes: 0,
			leadingOnes: 128, trailingOnes: 128,
		},
		{
			name:          "128 bits",
			s:             "01111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 1, trailingZeroes: 0,
			leadingOnes: 0, trailingOnes: 127,
		},
		{
			name:          "128 bits",
			s:             "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111110",
			leadingZeroes: 0, trailingZeroes: 1,
			leadingOnes: 127, trailingOnes: 0,
		},
		{
			name:          "128 bits",
			s:             "00000000000000000000000000000000000000000000000000000000000000001111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 64, trailingZeroes: 0,
			leadingOnes: 0, trailingOnes: 64,
		},
		{
			name:          "128 bits",
			s:             "11111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 0, trailingZeroes: 0,
			leadingOnes: 128, trailingOnes: 128,
		},
		{
			name:          "128 bits",
			s:             "00000000000000000000000000000000000000000000000000000000000000001111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 64, trailingZeroes: 0,
			leadingOnes: 0, trailingOnes: 64,
		},
		{
			name:          "128 bits",
			s:             "11111111111111111111111111111111111111111111111111111111111111110000000000000000000000000000000000000000000000000000000000000000",
			leadingZeroes: 0, trailingZeroes: 64,
			leadingOnes: 64, trailingOnes: 0,
		},

		{
			name:          "129 bits",
			s:             "111111111111111111111111111111111111111111111111111111111111111100000000000000000000000000000000000000000000000000000000000000000",
			leadingZeroes: 0, trailingZeroes: 65,
			leadingOnes: 64, trailingOnes: 0,
		},
		{
			name:          "129 bits",
			s:             "000000000000000000000000000000000000000000000000000000000000000001111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 65, trailingZeroes: 0,
			leadingOnes: 0, trailingOnes: 64,
		},
		{
			name:          "129 bits",
			s:             "000000000000000000000000000000000000000000000000000000000000000011111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 64, trailingZeroes: 0,
			leadingOnes: 0, trailingOnes: 65,
		},
		{
			name:          "129 bits",
			s:             "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000111000",
			leadingZeroes: 123, trailingZeroes: 3,
			leadingOnes: 0, trailingOnes: 0,
		},
		{
			name:          "129 bits",
			s:             "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111000111",
			leadingZeroes: 0, trailingZeroes: 0,
			leadingOnes: 123, trailingOnes: 3,
		},
		{
			name:          "129 bits",
			s:             "000111000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			leadingZeroes: 3, trailingZeroes: 123,
			leadingOnes: 0, trailingOnes: 0,
		},
		{
			name:          "129 bits",
			s:             "111000111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 0, trailingZeroes: 0,
			leadingOnes: 3, trailingOnes: 123,
		},
		{
			name:          "129 bits",
			s:             "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 0, trailingZeroes: 0,
			leadingOnes: 129, trailingOnes: 129,
		},
		{
			name:          "259 bits",
			s:             "0010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			leadingZeroes: 2, trailingZeroes: 136,
			leadingOnes: 0, trailingOnes: 0,
		},
		{
			name:          "259 bits",
			s:             "1101111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111101111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			leadingZeroes: 0, trailingZeroes: 0,
			leadingOnes: 2, trailingOnes: 136,
		},
	}

	for _, tt := range tests {
		bs, _ := NewFromString(tt.s)
		t.Run(tt.name, func(t *testing.T) {
			if got := bs.LeadingZeroes(); got != tt.leadingZeroes {
				t.Errorf("%q leading zeroes = %d, want %d", tt.s, got, tt.leadingZeroes)
			}

			if got := bs.TrailingZeroes(); got != tt.trailingZeroes {
				t.Errorf("%q trailing zeroes = %d, want %d", tt.s, got, tt.trailingZeroes)
			}

			if got := bs.LeadingOnes(); got != tt.leadingOnes {
				t.Errorf("%q leading ones = %d, want %d", tt.s, got, tt.leadingOnes)
			}

			if got := bs.TrailingOnes(); got != tt.trailingOnes {
				t.Errorf("%q trailing ones = %d, want %d", tt.s, got, tt.trailingOnes)
			}

			// After reversing the bit string, leading and trailing zeroes must be swapped.
			bs.Reverse()

			if got := bs.LeadingZeroes(); got != tt.trailingZeroes {
				t.Errorf("reversed %q leading zeroes = %d, want %d", tt.s, got, tt.trailingZeroes)
			}

			if got := bs.TrailingZeroes(); got != tt.leadingZeroes {
				t.Errorf("reversed  %q trailing zeroes = %d, want %d", tt.s, got, tt.leadingZeroes)
			}

			if got := bs.LeadingOnes(); got != tt.trailingOnes {
				t.Errorf("reversed %q leading ones = %d, want %d", tt.s, got, tt.trailingOnes)
			}

			if got := bs.TrailingOnes(); got != tt.leadingOnes {
				t.Errorf("reversed  %q trailing ones = %d, want %d", tt.s, got, tt.leadingOnes)
			}
		})
	}
}
