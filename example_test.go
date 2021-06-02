package bitstring

import (
	"fmt"
	"math/big"
)

func ExampleNew() {
	bs := New(32)

	fmt.Println(bs)
	// Output: 00000000000000000000000000000000
}

func ExampleNewFromString() {
	// Create a Bitstring from a string made of 0's and 1's.
	bs, _ := NewFromString("101001")

	fmt.Println(bs)
	fmt.Println(bs.Len(), "bits")
	// Output: 101001
	// 6 bits
}

func ExampleBitstring_SetBit() {
	bs := New(8)
	bs.SetBit(2)

	fmt.Println("bit 2:", bs.Bit(2))
	fmt.Println("bit 7:", bs.Bit(7))
	// Output: bit 2: true
	// bit 7: false
}

func ExampleBitstring_ClearBit() {
	bs, err := NewFromString("11111111")
	_ = err

	bs.ClearBit(2)

	fmt.Println(bs)
	// Output: 11111011
}

func ExampleBitstring_FlipBit() {
	bs := New(8)
	bs.FlipBit(2)

	fmt.Println(bs)
	// Output: 00000100
}

func ExampleBitstring_ZeroesCount() {
	bs := New(8)
	bs.SetBit(3)

	fmt.Println(bs.ZeroesCount())
	// Output: 7
}

func ExampleBitstring_OnesCount() {
	bs := New(8)
	bs.SetBit(3)

	fmt.Println(bs.OnesCount())
	// Output: 1
}

func ExampleNewFromBig() {
	var bi big.Int
	bi.SetString("11110000111100001111000011110000", 2)

	bs := NewFromBig(&bi)

	fmt.Println(bs)
	// Output: 11110000111100001111000011110000
}

func ExampleBitstring_BigInt() {
	bs, _ := NewFromString("100000")

	big := bs.BigInt().Int64()
	fmt.Println(big)
	// Output: 32
}

func ExampleSwapRange() {
	bs1, _ := NewFromString("11111001")
	bs2, _ := NewFromString("00000110")

	// Swap 2 bits from index 0
	SwapRange(bs1, bs2, 1, 2)

	fmt.Println(bs1)
	fmt.Println(bs2)
	// Output: 11111111
	// 00000000
}

func ExampleBitstring_ClearRange() {
	bs, _ := NewFromString("11111111")

	// Clear the 3 bits at offset 2.
	bs.ClearRange(2, 3)

	fmt.Println(bs)
	// Output: 11100011
}

func ExampleBitstring_SetRange() {
	bs, _ := NewFromString("00000000")

	// Set the 3 bits at offset 2.
	bs.SetRange(2, 3)

	fmt.Println(bs)
	// Output: 00011100
}

func ExampleBitstring_FlipRange() {
	bs, _ := NewFromString("11110000")

	// Flip the 3 bits at offset 2.
	bs.FlipRange(2, 4)

	fmt.Println(bs)
	// Output: 11001100
}

func ExampleBitstring_Flip() {
	bs, _ := NewFromString("11110000")

	bs.Flip()

	fmt.Println(bs)
	// Output: 00001111
}
