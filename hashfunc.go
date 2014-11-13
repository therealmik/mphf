package mphf

import (
	"github.com/dchest/siphash"
	"math/rand"
)

const (
	NUM_HASHES  = 3
)

type HashFunc struct {
	K0 uint64
	K1 uint64
	M  uint64
}

func RandomHashFunc(numBuckets uint) HashFunc {
	var k1 uint64 = (uint64(rand.Uint32()) << 32) | uint64(rand.Uint32())
	var k2 uint64 = (uint64(rand.Uint32()) << 32) | uint64(rand.Uint32())
	return HashFunc{k1, k2, uint64(numBuckets)}
}

func (self *HashFunc) Sum(data []byte) (ret [NUM_HASHES]uint) {
	h := siphash.Hash(self.K0, self.K1, data)
	ret[0] = uint((h >> 32) % self.M)
	ret[1] = uint(h % self.M)
	ret[2] = uint(((h >> 32) + h) % self.M)
	return ret
}

// Note: assumes M < (1 << 42)
func (self *HashFunc) BigSum(data []byte) (ret [NUM_HASHES]uint64) {
	h0, h1 := siphash.Hash128(self.K0, self.K1, data)
	ret[0] = h0 % self.M
	ret[1] = h1 % self.M

	upper0, upper1 := h0 >> 42, h1 >> 42
	ret[2] = (upper0 | (upper1 << 22)) % self.M

	return ret
}

