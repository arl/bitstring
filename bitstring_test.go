package bitstring

import (
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	bs := New(100)

	if bs.Len() != 100 {
		t.Errorf("len = %d, want 100", bs.Len())
	}

	for i := 0; i < bs.Len(); i++ {
		if bs.Bit(i) {
			t.Error("bit i = 1, want 0")
		}
	}
}

func TestRandom(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	bs := Random(100, rng)
	if bs.Len() != 100 {
		t.Errorf("len = %d, want 100", bs.Len())
	}
}

func TestSetBit(t *testing.T) {
	bs := New(69)

	bits := make([]bool, 69)

	bs.SetBit(1)
	bs.SetBit(4)
	bs.SetBit(65)
	bits[1] = true
	bits[4] = true
	bits[65] = true
	checkBits(t, bits, bs)

	bs.ClearBit(65)
	bits[65] = false
	checkBits(t, bits, bs)
}

func TestFlipBit(t *testing.T) {
	bs := New(69)

	bits := make([]bool, 69)

	bs.FlipBit(2)
	bits[2] = !bits[2]
	checkBits(t, bits, bs)

	bs.FlipBit(2)
	bits[2] = !bits[2]
	checkBits(t, bits, bs)

	bs.FlipBit(67)
	bits[67] = !bits[67]
	checkBits(t, bits, bs)

	bs.FlipBit(67)
	bits[67] = !bits[67]
	checkBits(t, bits, bs)
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
			bs, err := NewFromString(tt.str)
			if !tt.valid {
				if bs != nil || err == nil {
					t.Errorf("invalid string, got (%s, %v) want (nil, error)", bs, err)
				}
				return
			}
			if bs == nil || err != nil || bs.String() != tt.str {
				t.Errorf("valid string, got (%s, %v) want (%s, nil)", bs, err, tt.str)
			}
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
		"000000000000000000111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111",
	}
	for _, str := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(str)
			ones := bs.OnesCount()

			bs.Reverse()

			revstr := string(reverseBytes([]byte(str)))
			if bs.String() != revstr {
				t.Errorf("reversed string = %s, want %s", bs.String(), revstr)
			}
			if bs.OnesCount() != ones {
				t.Errorf("reversed string got %d ones, want %d", bs.OnesCount(), ones)
			}
		})
	}
}

func TestOnesCount(t *testing.T) {
	bs := New(65)
	if bs.OnesCount() != 0 {
		t.Fatalf("default constructed Bitstring shouldn't have any bit set, got %d", bs.OnesCount())
	}

	bs.SetBit(0)
	bs.SetBit(31)
	bs.SetBit(32)
	bs.SetBit(33)
	bs.SetBit(63)
	bs.SetBit(64)

	want := 6
	if got := bs.OnesCount(); got != want {
		t.Errorf("%s got %d ones, want %d", bs, got, want)
	}
}

func TestZeroesCount(t *testing.T) {
	bs := New(12)
	want := 12
	if got := bs.ZeroesCount(); got != want {
		t.Errorf("%s got %d zeroes, want %d", bs, got, want)
	}

	bs.SetBit(0)
	bs.SetBit(5)
	bs.SetBit(6)
	bs.SetBit(9)
	bs.SetBit(10)

	want = 7
	if got := bs.ZeroesCount(); got != want {
		t.Errorf("%s got %d zeroes, want %d", bs, got, want)
	}
}

func TestClone(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	bs := Random(2000, rng)
	bs.SetBit(3)
	bs.SetBit(7)
	bs.SetBit(8)

	equalbits(t, bs.Clone(), bs)
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

	if !org.Equals(org) {
		t.Errorf("org != org, want org == org")
	}
	if org.Equals(nil) {
		t.Errorf("org == nil, want org != nil")
	}

	if org.Equals(&Bitstring{}) {
		t.Errorf("org == default constructed Bitstring, want !=")
	}

	clone := org.Clone()
	if !clone.Equals(org) {
		t.Errorf("different bitstrings, clone=%s, org=%s", clone, org)
	}

	clone.FlipBit(0)
	if clone.Equals(org) {
		t.Errorf("same bitstrings, clone=%s, org=%s", clone, org)
	}

	// Bitstrings of different lengths but with the same bits set should not be equal.
	bs2 := New(9)
	bs2.SetBit(2)
	bs2.SetBit(5)
	bs2.SetBit(8)

	if bs2.Equals(org) {
		t.Errorf("same bitstrings, bs2=%s, org=%s", bs2, org)
	}
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
			want := trimLeadingZeroes(num)
			if sbig != want {
				t.Errorf("got big = %s, want %s", sbig, want)
			}

			// Create a bitstring from the binary string representation, create
			// a Bitstring from it and compare the significant bits.
			bs2 := NewFromBig(bi)

			// Remove leading zeroes since a big int doesn't have leading zeroes.
			got := bs.String()[bs.LeadingZeroes():]
			want = bs2.String()
			if got != want {
				t.Errorf("got big = %s, want %s", got, want)
			}
		})
	}
}
func TestLeadingTrailing(t *testing.T) {
	maxbits := 193

	for nbits := 1; nbits < maxbits; nbits++ {
		for i := 0; i < nbits; i++ {
			bs := New(nbits)

			t.Run("only zeroes/"+bs.String(), func(t *testing.T) {
				if bs.LeadingZeroes() != nbits {
					t.Fatalf("%v leading zeroes = %d, want %d", bs, bs.LeadingZeroes(), nbits)
				}
				if bs.TrailingZeroes() != nbits {
					t.Fatalf("%v trailing zeroes = %d, want %d", bs, bs.TrailingZeroes(), nbits)
				}
			})

			bs.Flip()

			t.Run("only ones/"+bs.String(), func(t *testing.T) {
				if bs.LeadingOnes() != nbits {
					t.Fatalf("%v leading ones = %d, want %d", bs, bs.LeadingOnes(), nbits)
				}
				if bs.TrailingOnes() != nbits {
					t.Fatalf("%v trailing ones = %d, want %d", bs, bs.TrailingOnes(), nbits)
				}
			})

			bs.Flip()

			bs.SetBit(i)

			t.Run(bs.String(), func(t *testing.T) {
				if bs.LeadingZeroes() != nbits-i-1 {
					t.Fatalf("%v leading zeroes = %d, want %d", bs, bs.LeadingZeroes(), nbits-i-1)
				}
				if bs.TrailingZeroes() != i {
					t.Fatalf("%v trailing zeroes = %d, want %d", bs, bs.TrailingZeroes(), i)
				}
			})

			bs.Flip()

			t.Run(bs.String(), func(t *testing.T) {
				if bs.LeadingOnes() != nbits-i-1 {
					t.Fatalf("%v leading ones = %d, want %d", bs, bs.LeadingOnes(), nbits-i-1)
				}
				if bs.TrailingOnes() != i {
					t.Fatalf("%v trailing ones = %d, want %d", bs, bs.TrailingOnes(), i)
				}
			})
		}
	}
}
