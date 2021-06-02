package bitstring

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lomask(t *testing.T) {
	tests := []struct {
		n    uint64
		want uint64
	}{
		{n: 0, want: 0b00000000000000000000000000000000},
		{n: 1, want: 0b00000000000000000000000000000001},
		{n: 2, want: 0b00000000000000000000000000000011},
		{n: 64 - 2, want: math.MaxUint64 >> 2},
		{n: 64 - 1, want: math.MaxUint64 >> 1},
		{n: 64, want: math.MaxUint64},
	}
	for _, tt := range tests {
		if got := lomask(tt.n); got != tt.want {
			t.Errorf("lomask(%d) got %s, want %s", tt.n,
				sprintubits(got, 32), sprintubits(tt.want, 32))
		}
	}
}

func Test_himask(t *testing.T) {
	tests := []struct {
		n    uint64
		want uint64
	}{
		{n: 0, want: math.MaxUint64},
		{n: 1, want: math.MaxUint64 - 1},
		{n: 64 - 2, want: 1<<(64-1) + 1<<(64-2)},
		{n: 64 - 1, want: 1 << (64 - 1)},
	}
	for _, tt := range tests {
		if got := himask(tt.n); got != tt.want {
			t.Errorf("himask(%d) got %s, want %s", tt.n,
				sprintubits(got, 32), sprintubits(tt.want, 32))
		}
	}
}

func Test_genmask(t *testing.T) {
	tests := []struct {
		l, h uint64
		want uint64
	}{
		{l: 0, h: 0, want: 0b00000000000000000000000000000000},
		{l: 0, h: 1, want: 0b00000000000000000000000000000001},
		{l: 0, h: 2, want: 0b00000000000000000000000000000011},
		{l: 1, h: 1, want: 0b00000000000000000000000000000000},
		{l: 1, h: 2, want: 0b00000000000000000000000000000010},
		{l: 0, h: 31, want: 0b01111111111111111111111111111111},
		{l: 1, h: 31, want: 0b01111111111111111111111111111110},
		{l: 0, h: 30, want: 0b00111111111111111111111111111111},
	}
	for _, tt := range tests {
		if got := mask(tt.l, tt.h); got != tt.want {
			t.Errorf("mask(%d, %d) got %02b, want %02b", tt.l, tt.h,
				got, tt.want)
		}
	}
}

func Test_lsb(t *testing.T) {
	tests := []struct {
		bits uint64
		want uint64
	}{
		{bits: 0b11000010000000000000000000000001, want: 0},
		{bits: 0b01000001000000000000000000000010, want: 1},
		{bits: 0b10000000000100000000000000000001, want: 0},
		{bits: 0b00000000010001111111000000000100, want: 2},
		{bits: 0b00000001000001111111000000000000, want: 12},
		{bits: 0b10000000000000000000000000000000, want: 31},
		{bits: 0b1000000000000000000000000000000000000000000000000000000000000000, want: 63},
		{bits: 0b11111111111111111111111111111111, want: 0},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := lsb(tt.bits)
			if tt.want != got {
				t.Errorf("lsb(%02b) = %d, want %d", tt.bits, got, tt.want)
			}
		})
	}
}

func Test_msb(t *testing.T) {
	tests := []struct {
		bits uint64
		want uint64
	}{
		{bits: 0b1100001000000000000000000000000100000000000000000000000000000000, want: 63},
		{bits: 0b01, want: 0},
		{bits: 0b10, want: 1},
		{bits: 0b10000000000100000000000000000001, want: 31},
		{bits: 0b110000000000100000000000000000001, want: 32},
		{bits: 0b00000000010001111111000000000100, want: 22},
		{bits: 0b00000001000001111111000000000000, want: 24},
		{bits: 0b01000000000000000000000000000000, want: 30},
		{bits: 0b1000000000000000000000000000000000000000000000000000000000000000, want: 63},
		{bits: 0b0100000000000000000000000000000000000000000000000000000000000000, want: 62},
		{bits: 0b1111111111111111111111111111111111111111111111111111111111111111, want: 63},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := msb(tt.bits)
			if tt.want != got {
				t.Errorf("msb(%02b) = %d, want %d", tt.bits, got, tt.want)
			}
		})
	}
}

func Test_ispow2(t *testing.T) {
	for i := uint64(0); i < 64; i++ {
		t.Run(fmt.Sprintf("1<<%d", i), func(t *testing.T) {
			n, ok := ispow2(1 << i)
			assert.Truef(t, ok, "ispow2(1<<%d)", i)
			assert.Equalf(t, i, n, "ispow2(1<<%d)", i)
			if i == 0 {
				return
			}
			_, ok = ispow2(1<<i + 1)
			assert.Falsef(t, ok, "ispow2(1<<%d)", i)
		})
	}
}
