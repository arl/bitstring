// +build ignore

package main

import (
	"fmt"
	"go/format"
	"log"
	"os"
	"strings"
)

// This program generates the lookup table used to reverse the bits in a byte,
// based on https://graphics.stanford.edu/~seander/bithacks.html#BitReverseTable

func main() {
	var sb strings.Builder

	r2 := func(n byte) []byte {
		return []byte{n, n + 2*64, n + 1*64, n + 3*64}
	}
	r4 := func(n byte) []byte {
		r := make([]byte, 0, 4*4)
		r = append(r, r2(n)...)
		r = append(r, r2(n+2*16)...)
		r = append(r, r2(n+1*16)...)
		r = append(r, r2(n+3*16)...)
		return r
	}
	r6 := func(n byte) []byte {
		r := make([]byte, 0, 4*4*4)
		r = append(r, r4(n)...)
		r = append(r, r4(n+2*4)...)
		r = append(r, r4(n+1*4)...)
		r = append(r, r4(n+3*4)...)
		return r
	}

	lut := make([]byte, 0, 256)
	lut = append(lut, r6(0)...)
	lut = append(lut, r6(2)...)
	lut = append(lut, r6(1)...)
	lut = append(lut, r6(3)...)

	sb.WriteString("// Generated code; DO NOT EDIT.\n")
	sb.WriteString("//\n")
	sb.WriteString("// generated with: go run reverse_lut_generate.go\n")
	sb.WriteString("\n")
	sb.WriteString("package bitstring\n")
	sb.WriteString("\n")
	sb.WriteString("// reverseLut maps a byte to the same byte with reversed bits.\n")
	sb.WriteString("var reverseLut = [256]byte {\n")

	for i := 0; i < 16; i++ {
		for j := 0; j < 16; j++ {
			fmt.Fprintf(&sb, "%d, ", lut[i*16+j])
		}
		sb.WriteByte('\n')
	}

	sb.WriteString("}\n")
	buf, err := format.Source([]byte(sb.String()))
	if err != nil {
		log.Fatal(err)
	}

	os.WriteFile("reverse_lut.go", buf, 0664)
}
