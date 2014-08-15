package mphf

import "math/big"

const (
	BITS_PER_COUNT = 255
)

type RankSelect struct {
	BitField *big.Int
	Counts   []byte
}

func NewRankSelect(numItems uint) RankSelect {
	numCounts := (numItems + BITS_PER_COUNT - 1) / BITS_PER_COUNT
	return RankSelect{big.NewInt(0), make([]byte, int(numCounts))}
}

func (self *RankSelect) Set(offset uint) {
	if self.BitField.Bit(int(offset)) == 0 {
		self.incCount(offset)
		self.BitField.SetBit(self.BitField, int(offset), 1)
	}
}

func (self *RankSelect) incCount(offset uint) {
	countNum := offset / BITS_PER_COUNT
	if len(self.Counts) <= int(countNum) {
		newCounts := make([]byte, int(countNum+1))
		copy(newCounts, self.Counts)
		self.Counts = newCounts
	}
	self.Counts[countNum]++
}

func (self RankSelect) Rank(offset uint) uint {
	countsToUse := offset / BITS_PER_COUNT

	var ret uint
	for i := uint(0); i < countsToUse && i < uint(len(self.Counts)); i++ {
		ret += uint(self.Counts[i])
	}

	startBit := countsToUse * BITS_PER_COUNT
	for i := startBit; i < offset; i++ {
		ret += self.BitField.Bit(int(i))
	}

	return ret
}
