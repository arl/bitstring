[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/arl/bitstring)
[![Test Actions Status](https://github.com/arl/bitstring/workflows/Test/badge.svg)](https://github.com/arl/bitstring/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/arl/bitstring)](https://goreportcard.com/report/github.com/arl/bitstring)
[![codecov](https://codecov.io/gh/arl/bitstring/branch/main/graph/badge.svg)](https://codecov.io/gh/arl/bitstring)

# `bitstring`

Go bitstring library

Package `bitstring` implements a fixed length bit string type and bit manipulation functions.

 - Get/Set/Clear/Flip a single bit: `Bit`|`SetBit`|`ClearBit`|`FlipBit`
 - Set/Clear/Flip a range of bits: `SetRange`|`ClearRange`|`FlipRange`
 - Compare 2 bit strings: `Equals` or `EqualsRange`
 - 8/16/32/64/N signed/unsigned to/from conversions:
   - `Uint8`|`Uint16`|`Uint32`|`Uint64`|`Uintn`
   - `SetUint8`|`SetUint16`|`SetUint32`|`SetUint64`|`SetUintn`
 - Count ones/zeroes: `ZeroesCount`|`OnesCount`
 - Gray code conversion methods: `Gray8`|`Gray16`|`Gray32`|`Gray64`|`Grayn`
 - Convert to/from `big.Int`: `BigInt` | `NewFromBig`
 - Copy/Clone methods: `Copy`|`Clone`|`CopyRange`
 - Trailing/Leading Zeroes/Ones : `TrailingZeroes`|`LeadingZeroes`|`TrailingOnes`|`LeadingOnes`


## Debug versionhttps://bitstring.readthedocs.io/en/latest/

By default, bit offsets arguments to `bitstring` methods are not checked. This
allows not to pay the performance penalty of always checking offsets, in
environments where they are constants or always known beforehand.

You can enable runtime checks by passing the `bitstring_debug` build tag to `go`
when building the `bitstring` package.

**TODO**:
 - RotateLeft/Right ShiftLeft/Right
 - Or, And, Xor between bitstrings

have a look at  https://bitstring.readthedocs.io/en/latest/
