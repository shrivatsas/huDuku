package roaring

import (
	"fmt"
)

const (
	bitmapContainerMaxCapacity = uint32(1 << 16)
	one                        = uint64(1)
)

type bitmapContainer struct {
	cardinality int
	bitmap      []uint64
}

func newBitmapContainer() *bitmapContainer {
	return &bitmapContainer{0, make([]uint64, bitmapContainerMaxCapacity/64)}
}

func (bc *bitmapContainer) loadData(ac *arrayContainer) {
	bc.cardinality = ac.cardinality
	for i := 0; i < ac.cardinality; i++ {
		bc.bitmap[uint32(ac.values[i])/64] |= one << (ac.values[i] % 64)
	}
}

func (bc *bitmapContainer) toArrayContainer() *arrayContainer {
	values := make([]uint16, bc.cardinality)
	pos := 0
	fmt.Println(bc.cardinality, len(bc.bitmap))
	for k := 0; k < len(bc.bitmap); k++ {
		bitset := bc.bitmap[k]
		for bitset != 0 {
			t := bitset & -bitset
			values[pos] = uint16((k*64 + countBits(t-1)))
			pos++
			bitset ^= t
		}
	}

	return &arrayContainer{bc.cardinality, values}
}

func (bc *bitmapContainer) add(i uint16) container {
	x := uint32(i)
	index := x / 64
	mod := x % 64
	previous := bc.bitmap[index]
	bc.bitmap[index] |= one << mod
	bc.cardinality += int((previous ^ bc.bitmap[index]) >> mod)
	return bc
}

func (bc *bitmapContainer) contains(x uint16) bool {
	return bc.bitmap[uint32(x)/64]&(one<<(x%64)) != 0
}

// http://en.wikipedia.org/wiki/Hamming_weight
func countBits(i uint64) int {
	i = i - ((i >> 1) & 0x5555555555555555)
	i = (i & 0x3333333333333333) + ((i >> 2) & 0x3333333333333333)
	result := (((i + (i >> 4)) & 0xF0F0F0F0F0F0F0F) * 0x101010101010101) >> 56
	return int(result)
}
