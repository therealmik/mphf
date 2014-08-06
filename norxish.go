package mphf

// Norxish - By Michael Samuel <mik@miknet.net>
//
// This is almost exactly the NORX AEAD from https://norx.io
//
// The changes I made are:
// - 2 rounds instead of 4 during data processing
// - No key (it's a hash) - data goes there instead
// - Allow a 64-bit personality - instead of 128-bit nonce
// - Return 3 64-bit integers - the first 3 of state after permutation
// - The length is *prepended* to the data
//
// The purpose of this is for constructing an MPHF - it is not intended
// for secure hashing, AEAD, or anything of the sort

const (
	R0 = 8
	R1 = 19
	R2 = 40
	R3 = 63
	U0 = 0x243f6a8885a308d3
	U1 = 0x13198a2e03707344
	U2 = 0xa4093822299f31d0
	U3 = 0x082efa98ec4e6c89
	U4 = 0xae8858dc339325a1
	U5 = 0x670a134ee52d7fa6
	U6 = 0xc4316d80cd967541
	U7 = 0xd21dfbf8b630b762
	U8 = 0x375a18d261e7f892
	U9 = 0x343d1f187d92285b
)

type Norxish [16]uint64

func rotr(x uint64, r uint) uint64 {
	return (x >> r) | (x << (64 - r))
}

func g(a, b, c, d uint64) (uint64, uint64, uint64, uint64) {
	a = (a ^ b) ^ ((a & b) << 1)
	d = rotr(a^d, R0)
	c = (c ^ d) ^ ((c & d) << 1)
	b = rotr(b^c, R1)

	a = (a ^ b) ^ ((a & b) << 1)
	d = rotr(a^d, R2)
	c = (c ^ d) ^ ((c & d) << 1)
	b = rotr(b^c, R3)

	return a, b, c, d
}

func norx_round(state *Norxish) {
	// Columns
	state[0], state[4], state[8], state[12] = g(state[0], state[4], state[8], state[12])
	state[1], state[5], state[9], state[13] = g(state[1], state[5], state[9], state[13])
	state[2], state[6], state[10], state[14] = g(state[2], state[6], state[10], state[14])
	state[3], state[7], state[11], state[15] = g(state[3], state[7], state[11], state[15])
	// Diagonals
	state[0], state[5], state[10], state[15] = g(state[0], state[5], state[10], state[15])
	state[1], state[6], state[11], state[12] = g(state[1], state[6], state[11], state[12])
	state[2], state[7], state[8], state[13] = g(state[2], state[7], state[8], state[13])
	state[3], state[4], state[9], state[14] = g(state[3], state[4], state[9], state[14])
}

func New(personality uint64) *Norxish {
	state := new(Norxish)

	state[0] = U0
	state[1] = personality
	// state[2] = 0
	state[3] = U1
	// state[3-7] = 0
	state[8] = U2
	state[9] = U3
	state[10] = U4
	state[11] = U5
	state[12] = U6
	state[13] = U7
	state[14] = U8
	state[15] = U9

	norx_round(state)
	norx_round(state)
	norx_round(state)
	norx_round(state)

	return state
}

func to_u64(data []byte) []uint64 {
	ret := make([]uint64, (len(data)+15)/8)
	ret[0] = uint64(len(data))

	for i := 0; i < len(data); i++ {
		var index int = (i / 8) + 1
		var shift uint = uint(i%8) * 8
		ret[index] |= uint64(data[i]) << shift
	}
	return ret
}

func min(x, y int) int {
	if x < y {
		return x
	} else {
		return y
	}
}

// Hash data, with personality and return 3 uint64 hash values.  WARNING: data > 32 bytes is dropped
func (s0 *Norxish) Hash(data []byte) (uint64, uint64, uint64) {
	var state *Norxish
	*state = *s0 // Copy the state

	input := to_u64(data)

	for len(input) > 0 {
		for i := 0; i < 8 && i < len(input); i++ {
			state[i] ^= input[i]
		}
		norx_round(state)
		norx_round(state)
		input = input[min(8, len(input)):]
	}
	return state[0], state[1], state[2]
}
