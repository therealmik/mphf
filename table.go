package mphf

import (
	"bytes"
)

const (
	INITIAL_TTL = 512
)

type KeyValue struct {
	Key   []byte
	Value interface{}
}

type TableEntry struct {
	Item    *KeyValue
	HashNum byte
}

type TableBuilder struct {
	Table []TableEntry
	Occup RankSelect
	Hash  HashFunc
	count uint
}

type MPHF struct {
	Table []KeyValue
	Occup RankSelect
	Hash  HashFunc
}

func NewTableBuilder(numItems uint) TableBuilder {
	numBuckets := (numItems * 5) / 4 // 80% Occup
	tbl := make([]TableEntry, numBuckets)
	occup := NewRankSelect(numItems)
	hash := RandomHashFunc(numBuckets)
	return TableBuilder{tbl, occup, hash, 0}
}

func (self *TableBuilder) displace(entry TableEntry, ttl int) bool {
	if ttl <= 0 {
		return false
	}

	bucket := self.Hash.Sum(entry.Item.Key)[int(entry.HashNum)%NUM_HASHES]
	displaced := self.Table[bucket]
	self.Table[bucket] = entry
	self.Occup.Set(bucket)

	if displaced.Item != nil {
		displaced.HashNum++
		return self.displace(displaced, ttl-1)
	} else {
		return true
	}
}

func (self *TableBuilder) Insert(item KeyValue) bool {
	entry := TableEntry{&item, 0}
	ret := self.displace(entry, INITIAL_TTL)
	if ret {
		self.count++
	}
	return ret
}

func (self *TableBuilder) MakeMPHF() MPHF {
	table := make([]KeyValue, 0, self.count)
	for _, x := range self.Table {
		if x.Item != nil {
			table = append(table, *x.Item)
		}
	}
	return MPHF{table, self.Occup, self.Hash}
}

func (self *MPHF) Get(key []byte) (interface{}, bool) {
	hs := self.Hash.Sum(key)
	for hn := 0; hn < NUM_HASHES; hn++ {
		phf_bucket := hs[hn]
		mphf_bucket := self.Occup.Rank(phf_bucket)
		item := self.Table[mphf_bucket]
		if bytes.Equal(item.Key, key) {
			return item.Value, true
		}
	}
	return nil, false
}

func BuildMPHF(items []KeyValue) (*MPHF, bool) {
	tb := NewTableBuilder(uint(len(items)))
	for _, x := range(items) {
		if !tb.Insert(x) {
			return nil, false
		}
	}
	mphf := tb.MakeMPHF()
	return &mphf, true
}
