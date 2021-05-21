package bitstring

import (
	"math"
	"testing"
)

func TestInt8(t *testing.T) {
	tests := []struct {
		input string
		i     int
		want  int8
	}{
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

func TestInt16(t *testing.T) {
	tests := []struct {
		input string
		i     int
		want  int16
	}{
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

func TestInt64(t *testing.T) {
	tests := []struct {
		input string
		i     int
		want  int64
	}{
		{
			i: 0, want: -1,
			input: "1111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			i: 0, want: math.MaxInt64,
			input: "0111111111111111111111111111111111111111111111111111111111111111",
		},
		{
			i: 0, want: math.MinInt64,
			input: "1000000000000000000000000000000000000000000000000000000000000000",
		},
		{
			i: 31, want: -1,
			input: "11111111111111111111111111111111111111111111111111111111111111110000000000000000000000000000000",
		},
		{
			i: 31, want: math.MaxInt64,
			input: "01111111111111111111111111111111111111111111111111111111111111110000000000000000000000000000000",
		},
		{
			i: 31, want: math.MinInt64,
			input: "10000000000000000000000000000000000000000000000000000000000000001111111111111111111111111111111",
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			bs, _ := NewFromString(tt.input)
			got := bs.Int64(tt.i)
			if tt.want != got {
				t.Errorf("Bitstring(%s).Int64(%d) got %s, want %s", tt.input, tt.i,
					sprintsbits(int64(got), 32), sprintsbits(int64(tt.want), 32))
			}
		})
	}
}

func TestSetInt(t *testing.T) {
	t.Run("SetInt8", func(t *testing.T) {
		bs := New(8)
		want := int8(math.MinInt8 / 3 * 2)
		bs.SetInt8(0, want)
		got := bs.Int8(0)
		if got != want {
			t.Errorf("Bitstring().SetInt8(%d, %d) got %d, want %d",
				0, want, got, want)
		}
	})
	t.Run("SetInt16", func(t *testing.T) {
		bs := New(16)
		want := int16(math.MinInt16 / 3 * 2)
		bs.SetInt16(0, want)
		got := bs.Int16(0)
		if got != want {
			t.Errorf("Bitstring().SetInt16(%d, %d) got %d, want %d",
				0, want, got, want)
		}
	})
	t.Run("SetInt32", func(t *testing.T) {
		bs := New(32)
		want := int32(math.MinInt32 / 3 * 2)
		bs.SetInt32(0, want)
		got := bs.Int32(0)
		if got != want {
			t.Errorf("Bitstring().SetInt32(%d, %d) got %d, want %d",
				0, want, got, want)
		}
	})
	t.Run("SetInt64", func(t *testing.T) {
		bs := New(64)
		want := int64(math.MinInt64 / 3 * 2)
		bs.SetInt64(0, want)
		got := bs.Int64(0)
		if got != want {
			t.Errorf("Bitstring().SetInt64(%d, %d) got %d, want %d",
				0, want, got, want)
		}
	})
	t.Run("SetIntn", func(t *testing.T) {
		bs := New(97)
		want := int64(math.MinInt64 / 3 * 2)
		bs.SetIntn(12, 45, want)
		got := bs.Intn(12, 45)
		if got != want&int64(lomask(45)) {
			t.Errorf("Bitstring().SetInt64(%d, %d) got %d, want %d",
				0, want, got, want)
		}
	})
}
