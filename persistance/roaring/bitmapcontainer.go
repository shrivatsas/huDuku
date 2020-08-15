package roaring

const (
	bitmapContainerMaxCapacity = uint32(1 << 16)
	// one                        = uint64(1)
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
		bc.bitmap[uint32(ac.content[i])/64] |= one << (ac.content[i] % 64)
	}
}

func (bc *bitmapContainer) toArrayContainer() *arrayContainer {
	values := make([]uint16, bc.cardinality)
	pos := 0
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
