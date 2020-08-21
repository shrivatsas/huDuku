package roaring

type entry struct {
	key       uint16
	container container
}

type Roaring struct {
	containers []entry
}

func bitmapToArray(bc *bitmapContainer) *arrayContainer {
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

func arrayToBitmap(ac *arrayContainer) *bitmapContainer {
	bc := newBitmapContainer()
	bc.cardinality = ac.cardinality
	for i := 0; i < ac.cardinality; i++ {
		bc.bitmap[uint32(ac.values[i])/64] |= one << (ac.values[i] % 64)
	}
	return bc
}

func and(c1 *bitmapContainer, c2 *arrayContainer) container {
	answer := make([]uint16, c2.cardinality)

	cardinality := 0
	for k := 0; k < c2.cardinality; k++ {
		if c1.contains(c2.values[k]) {
			answer[cardinality] = c2.values[k]
			cardinality++
		}
	}

	return &arrayContainer{cardinality, answer[:cardinality]}
}

func or(c1 *bitmapContainer, c2 *arrayContainer) container {
	bitmap := make([]uint64, c1.cardinality)
	copy(bitmap, c1.bitmap)

	answer := &bitmapContainer{c1.cardinality, bitmap}
	for _, v := range c2.values {
		answer.add(v)
	}

	return answer
}
