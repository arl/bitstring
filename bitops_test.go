package bitstring

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_lomask(t *testing.T) {
	tests := []struct {
		n    uint64
		want uint64
	}{
		{n: 0, want: atobin("00000000000000000000000000000000")},
		{n: 1, want: atobin("00000000000000000000000000000001")},
		{n: 2, want: atobin("00000000000000000000000000000011")},
		{n: uintsize - 2, want: maxuint >> 2},
		{n: uintsize - 1, want: maxuint >> 1},
		{n: uintsize, want: maxuint},
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
		{n: 0, want: maxuint},
		{n: 1, want: maxuint - 1},
		{n: uintsize - 2, want: 1<<(uintsize-1) + 1<<(uintsize-2)},
		{n: uintsize - 1, want: 1 << (uintsize - 1)},
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
		want string
	}{
		{l: 0, h: 0, want: "00000000000000000000000000000000"},
		{l: 0, h: 1, want: "00000000000000000000000000000001"},
		{l: 0, h: 2, want: "00000000000000000000000000000011"},
		{l: 1, h: 1, want: "00000000000000000000000000000000"},
		{l: 1, h: 2, want: "00000000000000000000000000000010"},
		{l: 0, h: 31, want: "01111111111111111111111111111111"},
		{l: 1, h: 31, want: "01111111111111111111111111111110"},
		{l: 0, h: 30, want: "00111111111111111111111111111111"},
	}
	for _, tt := range tests {
		if got := mask(tt.l, tt.h); got != atobin(tt.want) {
			t.Errorf("mask(%d, %d) got %s, want %s", tt.l, tt.h,
				sprintubits(got, 32), tt.want)
		}
	}
}

func Test_firstSetBit(t *testing.T) {
	tests := []struct {
		bits string
		want uint64
	}{
		{bits: "00000000000000000000000000000001", want: 0},
		{bits: "00000000000000000000000000000010", want: 1},
		{bits: "10000000000000000000000000000001", want: 0},
		{bits: "00000000000001111111000000000100", want: 2},
		{bits: "00000000000001111111000000000000", want: 12},
		{bits: "10000000000000000000000000000000", want: 31},
		{bits: "00000000000000000000000000000000", want: uintsize - 1},
		{bits: "1000000000000000000000000000000000000000000000000000000000000000", want: 63},
		{bits: "11111111111111111111111111111111", want: 0},
	}
	for _, tt := range tests {
		got := firstSetBit(atobin(tt.bits))
		assert.EqualValuesf(t, tt.want, got, "%q, first bit = %d, want %d", tt.bits, got, tt.want)
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
