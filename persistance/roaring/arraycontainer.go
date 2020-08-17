package roaring

const (
	arrayContainerInitSize = 4
	arrayContainerMaxSize  = 4096
)

type arrayContainer struct {
	cardinality int
	values      []uint16
}

func newArrayContainer() *arrayContainer {
	values := make([]uint16, arrayContainerInitSize)
	return &arrayContainer{0, values}
}

func newArrayContainerWithCapacity(capacity int) *arrayContainer {
	values := make([]uint16, capacity)
	return &arrayContainer{0, values}
}

func newArrayContainerRunOfOnes(firstOfRun, lastOfRun int) *arrayContainer {
	valuesInRange := lastOfRun - firstOfRun + 1
	content := make([]uint16, valuesInRange)
	for i := 0; i < valuesInRange; i++ {
		content[i] = uint16(firstOfRun + i)
	}
	return &arrayContainer{int(valuesInRange), content}
}

func (ac *arrayContainer) contains(x uint16) bool {
	return binarySearch(ac.values, ac.cardinality, x) >= 0
}

func (ac *arrayContainer) add(x uint16) container {
	if ac.cardinality >= arrayContainerMaxSize {
		bc := ac.toBitmapContainer()
		bc.add(x)
		return bc
	}

	if ac.cardinality == 0 || x > ac.values[ac.cardinality-1] {
		if ac.cardinality >= len(ac.values) {
			ac.increaseCapacity()
		}
		ac.values[ac.cardinality] = x
		ac.cardinality++
		return ac
	}

	loc := binarySearch(ac.values, ac.cardinality, x)
	if loc < 0 {
		if ac.cardinality >= len(ac.values) {
			ac.increaseCapacity()
		}
		loc = -loc - 1
		// insertion : shift the elements > x by one position to
		// the right and put x in it's appropriate place
		copy(ac.values[loc+1:], ac.values[loc:])
		ac.values[loc] = x
		ac.cardinality++
	}
	return ac
}

func (ac *arrayContainer) increaseCapacity() {
	length := len(ac.values)
	var newLength int
	switch {
	case length < 64:
		newLength = length * 2
	case length < 1024:
		newLength = length * 3 / 2
	default:
		newLength = length * 5 / 4
	}
	if newLength > arrayContainerMaxSize {
		newLength = arrayContainerMaxSize
	}
	newSlice := make([]uint16, newLength)
	copy(newSlice, ac.values)
	ac.values = newSlice
}

func (ac *arrayContainer) toBitmapContainer() *bitmapContainer {
	bc := newBitmapContainer()
	bc.loadData(ac)
	return bc
}

func binarySearch(array []uint16, length int, k uint16) int {
	low := 0
	high := length - 1

	for low <= high {
		middleIndex := (low + high) >> 1
		middleValue := array[middleIndex]

		switch {
		case middleValue < k:
			low = middleIndex + 1
		case middleValue > k:
			high = middleIndex - 1
		default:
			return middleIndex
		}
	}
	return -(low + 1)
}
