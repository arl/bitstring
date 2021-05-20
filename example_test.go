package bitstring

import "fmt"

func ExampleNew() {
	// Create a 32-bit Bitstring
	bs := New(32)
	// upon creation all bits are unset
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
	bs := New(8)
	bs.SetBit(2)
	fmt.Println(bs)
	bs.ClearBit(2)
	fmt.Println(bs)
	// Output: 00000100
	// 00000000
}

func ExampleBitstring_FlipBit() {
	bs := New(8)
	bs.FlipBit(2)
	fmt.Println(bs)
	// Output: 00000100
}

func ExampleBitstring_ZeroesCount() {
	bs := New(8)
	fmt.Println(bs.ZeroesCount())
	// Output: 8
}

func ExampleBitstring_OnesCount() {
	bitstring := New(8)
	fmt.Println(bitstring.OnesCount())
	// Output: 0
}

func ExampleBitstring_BigInt() {
	bs, _ := NewFromString("100")
	bi := bs.BigInt()
	fmt.Println(bi.Int64())
	// Output: 4
}

func ExampleSwapRange() {
	bs1, _ := NewFromString("111")
	bs2, _ := NewFromString("000")
	// Swap 2 bits from index 0
	SwapRange(bs1, bs2, 2, 1)
	fmt.Println(bs1)
	// Output: 011
}

func ExampleBitstring_ClearRange() {
	bs, _ := NewFromString("10101010")
	// Clear the 3 bits at offset 2.
	bs.ClearRange(2, 3)
	fmt.Println(bs)
	// Output: 10100010
}

func ExampleBitstring_SetRange() {
	bs, _ := NewFromString("10101010")
	// Set the 3 bits at offset 2.
	bs.SetRange(2, 3)
	fmt.Println(bs)
	// Output: 10111110
}

func ExampleBitstring_FlipRange() {
	bs, _ := NewFromString("10101010")
	// Flip the 3 bits at offset 2.
	bs.FlipRange(2, 3)
	fmt.Println(bs)
	// Output: 10110110
}
