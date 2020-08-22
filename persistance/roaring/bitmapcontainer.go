package roaring

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

func (bc *bitmapContainer) add(i uint16) container {
	x := uint32(i)
	index := x / 64
	mod := x % 64
	previous := bc.bitmap[index]
	bc.bitmap[index] |= one << mod
	bc.cardinality += int((previous ^ bc.bitmap[index]) >> mod)
	return bc
}

func (bc *bitmapContainer) and(c1 container) container {
	switch oc := c1.(type) {
	case *bitmapContainer:
		newCardinality := 0
		for k, v := range bc.bitmap {
			newCardinality += countBits(v & oc.bitmap[k])
		}

		if newCardinality > arrayContainerMaxSize {
			answer := newBitmapContainer()
			for k, v := range bc.bitmap {
				answer.bitmap[k] = v & oc.bitmap[k]
			}
			answer.cardinality = newCardinality
			return answer
		}
		content := fillArrayAND(bc.bitmap, oc.bitmap, newCardinality)
		return &arrayContainer{newCardinality, content}
	case *arrayContainer:
		return and(bc, oc)
	}

	return nil
}

func (bc *bitmapContainer) or(c1 container) container {
	switch oc := c1.(type) {
	case *bitmapContainer:
		answer := newBitmapContainer()
		for i := 0; i < len(bc.bitmap); i++ {
			answer.bitmap[i] = bc.bitmap[i] | oc.bitmap[i]
			answer.cardinality = answer.cardinality + countBits(answer.bitmap[i])
		}
		if answer.getCardinality() < arrayContainerMaxSize {
			return bitmapToArray(answer)
		}
		return answer
	case *arrayContainer:
		return or(bc, oc)
	}
	return nil
}

func (bc *bitmapContainer) contains(x uint16) bool {
	return bc.bitmap[uint32(x)/64]&(one<<(x%64)) != 0
}

func (bc *bitmapContainer) getCardinality() int {
	return bc.cardinality
}

// http://en.wikipedia.org/wiki/Hamming_weight :: popcount64c
func countBits(i uint64) int {
	i = i - ((i >> 1) & 0x5555555555555555)
	i = (i & 0x3333333333333333) + ((i >> 2) & 0x3333333333333333)
	result := (((i + (i >> 4)) & 0xF0F0F0F0F0F0F0F) * 0x101010101010101) >> 56
	return int(result)
}

func fillArrayAND(bitmap1, bitmap2 []uint64, newCardinality int) []uint16 {
	pos := 0

	if len(bitmap1) != len(bitmap2) {
		panic("Bitmaps have different length - not supported.")
	}

	container := make([]uint16, newCardinality)
	for k := 0; k < len(bitmap1); k++ {
		bitset := bitmap1[k] & bitmap2[k]
		for bitset != 0 {
			t := bitset & -bitset
			container[pos] = uint16((k*64 + countBits(t-1)))
			pos++
			bitset ^= t
		}
	}

	return container
}
