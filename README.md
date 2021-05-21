[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/arl/bitstring)
[![Test Actions Status](https://github.com/arl/bitstring/workflows/Test/badge.svg)](https://github.com/arl/bitstring/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/arl/bitstring)](https://goreportcard.com/report/github.com/arl/bitstring)
[![codecov](https://codecov.io/gh/arl/bitstring/branch/main/graph/badge.svg)](https://codecov.io/gh/arl/bitstring)

# bitstring
Go bitstring library

Package `bitstring` implements a fixed-length bit string type and many bit manipulation functions:
 - set/clear/flip a single bit 
 - set/clear/flip a range of bits 
 - swap or compare range of bits between 2 bitstrings
 - 8/16/32/64/n signed/unsigned to/from conversions
 - count ones/zeroes
 - gray code conversion methods
 - convert to `big.Int`
 - Copy/Clone methods

TODO:
 - RotateLeft/Right
 - Trailing/Leading zeroes/ones
 - improve documentation
