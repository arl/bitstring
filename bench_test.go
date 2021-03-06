package bitstring

import (
	"math/rand"
	"testing"
)

var sink interface{}

func benchmarkUintn(b *testing.B, nbits, i int) {
	b.ReportAllocs()
	bs, _ := NewFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v uint64
	for n := 0; n < b.N; n++ {
		v = bs.Uintn(i, nbits)
	}
	b.StopTimer()
	sink = v
}

func benchmarkUint64(b *testing.B, i int) {
	b.ReportAllocs()
	bs, _ := NewFromString("00000000000000000000000000000001010000000000000000000000000000000")
	var v uint64
	for n := 0; n < b.N; n++ {
		v = bs.Uint64(i)
	}
	b.StopTimer()
	sink = v
}

func benchmarkUint32(b *testing.B, i int) {
	b.ReportAllocs()
	bs, _ := NewFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v uint32
	for n := 0; n < b.N; n++ {
		v = bs.Uint32(i)
	}
	b.StopTimer()
	sink = v
}

func benchmarkUint16(b *testing.B, i int) {
	b.ReportAllocs()
	bs, _ := NewFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v uint16
	for n := 0; n < b.N; n++ {
		v = bs.Uint16(i)
	}
	b.StopTimer()
	sink = v
}

func benchmarkUint8(b *testing.B, i int) {
	b.ReportAllocs()
	bs, _ := NewFromString("0000000000000000000000000000000101000000000000000000000000000000")
	var v uint8
	for n := 0; n < b.N; n++ {
		v = bs.Uint8(i)
	}
	b.StopTimer()
	sink = v
}

func BenchmarkUintnSameWord(b *testing.B)        { benchmarkUintn(b, 32, 32) }
func BenchmarkUintnDifferentWords(b *testing.B)  { benchmarkUintn(b, 32, 31) }
func BenchmarkUint64SameWord(b *testing.B)       { benchmarkUint64(b, 0) }
func BenchmarkUint64DifferentWords(b *testing.B) { benchmarkUint64(b, 1) }
func BenchmarkUint32SameWord(b *testing.B)       { benchmarkUint32(b, 32) }
func BenchmarkUint32DifferentWords(b *testing.B) { benchmarkUint32(b, 31) }
func BenchmarkUint16SameWord(b *testing.B)       { benchmarkUint16(b, 32) }
func BenchmarkUint16DifferentWords(b *testing.B) { benchmarkUint16(b, 31) }
func BenchmarkUint8SameWord(b *testing.B)        { benchmarkUint8(b, 32) }
func BenchmarkUint8DifferentWords(b *testing.B)  { benchmarkUint8(b, 31) }

func Benchmark_mask(b *testing.B) {
	b.ReportAllocs()

	var v uint64
	for i := 0; i < b.N; i++ {
		v = mask(4, 27)
	}
	b.StopTimer()
	sink = v
}

func Benchmark_lomask(b *testing.B) {
	b.ReportAllocs()

	var v uint64
	for i := 0; i < b.N; i++ {
		v = lomask(27)
	}
	b.StopTimer()
	sink = v
}

func BenchmarkSwapRange(b *testing.B) {
	x := New(1026)
	y := New(1026)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		SwapRange(x, y, i%26, 1026)
	}
}

func BenchmarkRandom(b *testing.B) {
	var x *Bitstring

	rng := rand.New(rand.NewSource(99))
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		x = Random(1026, rng)
	}
	b.StopTimer()
	sink = x
}

func benchmarkEquals(len int) func(b *testing.B) {
	return func(b *testing.B) {
		x, y := New(len), New(len)
		var res bool
		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			res = x.Equals(y)
		}
		b.StopTimer()
		sink = res
	}
}

func BenchmarkEquals(b *testing.B) {
	b.Run("len=100", benchmarkEquals(100))
	b.Run("len=500", benchmarkEquals(500))
	b.Run("len=1000", benchmarkEquals(1000))
	b.Run("len=2000", benchmarkEquals(2000))
	b.Run("len=4000", benchmarkEquals(4000))
	b.Run("len=8000", benchmarkEquals(8000))
}

func BenchmarkSetUint8(b *testing.B) {
	bs := New(67)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		bs.SetUint8(59, 255)
	}
	b.StopTimer()
	sink = bs
}

func BenchmarkSetUintn(b *testing.B) {
	bs := New(117)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		bs.SetUintn(35, 64, 0x9cfbeb71ee3fcf5f&64)
	}
	b.StopTimer()
	sink = bs
}

var vals = []string{
	"1111111111111111",
	"00000000000000000000000000000000000000000",
	"1111111111111111111111111111111111111111111111111111",
	"11111111111111111111111111111111111111111111111111111111100100111111100110010101010000010100111",
	"1111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111111100100111111100110010101010000010100111",
}

func benchmarkOnesCount(b *testing.B, val string) {
	bs, err := NewFromString(val)
	if err != nil {
		b.Fatal(err)
	}

	var ones int
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		ones = bs.OnesCount()
	}
	b.StopTimer()
	sink = ones
}

func BenchmarkOnesCount(b *testing.B) {
	for _, tt := range vals {
		b.Run("", func(b *testing.B) { benchmarkOnesCount(b, tt) })
	}
}

func Benchmark_msb(b *testing.B) {
	nums := []uint64{
		atobin("1100001000000000000000000000000100000000000000000000000000000000"),
		atobin("01"),
		atobin("10"),
		atobin("10000000000100000000000000000001"),
		atobin("110000000000100000000000000000001"),
		atobin("00000000010001111111000000000100"),
		atobin("00000001000001111111000000000000"),
		atobin("01000000000000000000000000000000"),
		atobin("1000000000000000000000000000000000000000000000000000000000000000"),
		atobin("0100000000000000000000000000000000000000000000000000000000000000"),
		atobin("1111111111111111111111111111111111111111111111111111111111111111"),
		atobin("1011111111111111111111111111111111111111111111111111111111111111"),
		atobin("1111111111111111111111111111111111111111111111111111111111111110"),
	}

	b.ResetTimer()
	var val uint64
	for i := 0; i < b.N; i++ {
		for _, num := range nums {
			val = msb(num)
		}
	}

	sink = val
}
