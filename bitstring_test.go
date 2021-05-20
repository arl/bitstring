package bitstring

import (
	"math/rand"
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
			got, err := MakeFromString(tt.str)
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

// Checks that integer conversion is correct.
func TestToNumber(t *testing.T) {
	bs := New(10)

	bs.SetBit(0)
	bs.SetBit(9)
	bint := bs.BigInt()
	assert.True(t, bint.IsInt64())
	assert.EqualValues(t, 513, bint.Int64())
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

	cpy := Clone(bs)
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

	clone := Clone(org)
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

func TestSwapRange(t *testing.T) {
	tests := []struct {
		x, y          string
		start, length int
		wantx, wanty  string
	}{
		{
			x:     "1",
			y:     "0",
			start: 0, length: 1,
			wantx: "0",
			wanty: "1",
		},
		{
			x:     "1111111111111111111111111111111111111111111111111111111111111111",
			y:     "0000000000000000000000000000000000000000000000000000000000000000",
			start: 0, length: 32,
			wantx: "1111111111111111111111111111111100000000000000000000000000000000",
			wanty: "0000000000000000000000000000000011111111111111111111111111111111",
		},
		{
			x:     "1111111111111111111111111111111111111111",
			y:     "0000000000000000000000000000000000000000",
			start: 2, length: 30,
			wantx: "1111111100000000000000000000000000000011",
			wanty: "0000000011111111111111111111111111111100",
		},
		{
			x:     "1111111111",
			y:     "0000000000",
			start: 0, length: 3,
			wantx: "1111111000",
			wanty: "0000000111",
		},
		{
			x:     "111",
			y:     "000",
			start: 1, length: 2,
			wantx: "001",
			wanty: "110",
		},
		{
			x:     "111",
			y:     "000",
			start: 0, length: 3,
			wantx: "000",
			wanty: "111",
		},
		{
			x:     "11111111111111111111111111111111",
			y:     "00000000000000000000000000000000",
			start: 0, length: 32,
			wantx: "00000000000000000000000000000000",
			wanty: "11111111111111111111111111111111",
		},
		{
			x:     "111111111111111111111111111111111",
			y:     "000000000000000000000000000000000",
			start: 0, length: 33,
			wantx: "000000000000000000000000000000000",
			wanty: "111111111111111111111111111111111",
		},
		{
			x:     "111111111111111111111111111111111111111111111111111111111111111",
			y:     "000000000000000000000000000000000000000000000000000000000000000",
			start: 0, length: 63,
			wantx: "000000000000000000000000000000000000000000000000000000000000000",
			wanty: "111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x:     "1111111111111111111111111111111111111111111111111111111111111111",
			y:     "0000000000000000000000000000000000000000000000000000000000000000",
			start: 0, length: 64,
			wantx: "0000000000000000000000000000000000000000000000000000000000000000",
			wanty: "1111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x:     "11111111111111111111111111111111111111111111111111111111111111111",
			y:     "00000000000000000000000000000000000000000000000000000000000000000",
			start: 0, length: 65,
			wantx: "00000000000000000000000000000000000000000000000000000000000000000",
			wanty: "11111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			x:     "1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			y:     "0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			start: 94, length: 1,
			wantx: "1101111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			wanty: "0010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			x:     "111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
			y:     "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			start: 1, length: 256,
			wantx: "100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
			wanty: "011111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111110",
		},
		{
			x:     "1111111111111111111111111111111111111111111111111111111111111111111",
			y:     "0000000000000000000000000000000000000000000000000000000000000000000",
			start: 64, length: 2,
			wantx: "1001111111111111111111111111111111111111111111111111111111111111111",
			wanty: "0110000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			x:     "1111111111111111111111111111111111111111111111111111111111111111111",
			y:     "0000000000000000000000000000000000000000000000000000000000000000000",
			start: 65, length: 1,
			wantx: "1011111111111111111111111111111111111111111111111111111111111111111",
			wanty: "0100000000000000000000000000000000000000000000000000000000000000000",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			x, err1 := MakeFromString(tt.x)
			assert.NoError(t, err1)
			y, err2 := MakeFromString(tt.y)
			assert.NoError(t, err2)
			SwapRange(x, y, tt.start, tt.length)

			assert.Equal(t, tt.wantx, x.String())
			assert.Equal(t, tt.wanty, y.String())
		})
	}
}
