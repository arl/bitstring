// +build !dev

package bitstring

// mustExist is a no-op in dev mode.
func (bs *Bitstring) mustExist(i uint) {}
