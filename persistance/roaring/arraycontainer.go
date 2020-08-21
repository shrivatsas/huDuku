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
		bc := arrayToBitmap(ac)
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

func (ac *arrayContainer) and(c1 container) container {
	switch oc := c1.(type) {
	case *arrayContainer:
		cardinality, content := intersect2by2(ac.values,
			ac.cardinality, oc.values,
			oc.getCardinality())
		return &arrayContainer{cardinality, content}
	case *bitmapContainer:
		return and(oc, ac)
	}

	return nil
}

func (ac *arrayContainer) or(c1 container) container {
	switch oc := c1.(type) {
	case *arrayContainer:
	case *bitmapContainer:
		return or(oc, ac)
	}

	return nil
}

func (ac *arrayContainer) orArray(other *arrayContainer) container {
	totalCardinality := ac.cardinality + other.cardinality
	if totalCardinality > arrayContainerMaxSize {
		bc := newBitmapContainer()
		for i := 0; i < other.cardinality; i++ {
			bc.add(other.values[i])
		}
		for i := 0; i < ac.cardinality; i++ {
			bc.add(ac.values[i])
		}
		if bc.cardinality <= arrayContainerMaxSize {
			return bitmapToArray(bc)
		}
		return bc
	}
	answer := arrayContainer{}
	pos, content := union2by2(ac.values, ac.cardinality, other.values, other.cardinality, totalCardinality)
	answer.cardinality = pos
	answer.values = content
	return &answer
}

func (ac *arrayContainer) andNotArray(value2 *arrayContainer) *arrayContainer {
	cardinality, content := difference(ac.values, ac.cardinality,
		value2.values, value2.cardinality)

	return &arrayContainer{cardinality, content}
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

func (ac *arrayContainer) getCardinality() int {
	return ac.cardinality
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

// Unite two sorted lists
func union2by2(set1 []uint16, length1 int,
	set2 []uint16, length2, bufferSize int) (int, []uint16) {

	if 0 == length2 {
		buffer := make([]uint16, length1)
		copy(buffer, set1)
		return length1, buffer
	}

	if 0 == length1 {
		buffer := make([]uint16, length2)
		copy(buffer, set2)
		return length2, buffer
	}

	buffer := make([]uint16, bufferSize)

	k1, k2, pos := 0, 0, 0

	for {
		if set1[k1] < set2[k2] {
			buffer[pos] = set1[k1]
			pos = pos + 1
			k1 = k1 + 1
			if k1 >= length1 {
				for ; k2 < length2; k2++ {
					buffer[pos] = set2[k2]
					pos = pos + 1
				}
				break
			}
		} else if set1[k1] == set2[k2] {
			buffer[pos] = set1[k1]
			pos = pos + 1
			k1 = k1 + 1
			k2 = k2 + 1
			if k1 >= length1 {
				for ; k2 < length2; k2++ {
					buffer[pos] = set2[k2]
					pos = pos + 1
				}
				break
			}
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos = pos + 1
				}
				break
			}
		} else {
			buffer[pos] = set2[k2]
			pos = pos + 1
			k2 = k2 + 1
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos = pos + 1
				}
				break
			}
		}
	}
	return pos, buffer[:pos]
}

func difference(
	set1 []uint16, length1 int,
	set2 []uint16, length2 int) (int, []uint16) {

	k1, k2, pos := 0, 0, 0

	if 0 == length2 {
		buffer := make([]uint16, length1)
		copy(buffer, set1)
		return length1, buffer
	}

	if 0 == length1 {
		return 0, make([]uint16, 0)
	}

	buffer := make([]uint16, length1)

	for {
		if set1[k1] < set2[k2] {
			buffer[pos] = set1[k1]
			pos++
			k1++
			if k1 >= length1 {
				break
			}
		} else if set1[k1] == set2[k2] {
			k1++
			k2++
			if k1 >= length1 {
				break
			}
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos++
				}
				break
			}
		} else { // if (val1>val2)
			k2++
			if k2 >= length2 {
				for ; k1 < length1; k1++ {
					buffer[pos] = set1[k1]
					pos++
				}
				break
			}
		}
	}
	return pos, buffer[:pos]
}

func intersect2by2(set1 []uint16, length1 int,
	set2 []uint16, length2 int) (int, []uint16) {

	if length1*64 < length2 {
		return oneSidedGallopingIntersect2by2(set1, length1, set2, length2)
	}

	if length2*64 < length1 {
		return oneSidedGallopingIntersect2by2(set2, length2, set1, length1)
	}

	return localIntersect2by2(set1, length1, set2, length2)
}

func min(x, y int) int {
	if x <= y {
		return x
	}
	return y
}

func localIntersect2by2(set1 []uint16, length1 int,
	set2 []uint16, length2 int) (int, []uint16) {

	if 0 == length1 || 0 == length2 {
		return 0, make([]uint16, 0)
	}

	finalLength := min(length1, length2)
	buffer := make([]uint16, finalLength)
	k1, k2, pos := 0, 0, 0

Mainwhile:
	for {
		if set2[k2] < set1[k1] {
			for {
				k2++
				if k2 == length2 {
					break Mainwhile
				}
				if set2[k2] >= set1[k1] {
					break
				}
			}
		}
		if set1[k1] < set2[k2] {
			for {
				k1++
				if k1 == length1 {
					break Mainwhile
				}
				if set1[k1] >= set2[k2] {
					break
				}
			}
		} else {
			buffer[pos] = set1[k1]
			pos++
			k1++
			if k1 == length1 {
				break
			}
			k2++
			if k2 == length2 {
				break
			}
		}
	}
	return pos, buffer[:pos]
}

func oneSidedGallopingIntersect2by2(
	smallSet []uint16, smallLength int,
	largeSet []uint16, largeLength int) (int, []uint16) {

	if 0 == smallLength {
		return 0, make([]uint16, 0)
	}

	buffer := make([]uint16, smallLength)
	k1, k2, pos := 0, 0, 0

	for {
		if largeSet[k1] < smallSet[k2] {
			k1 = advanceUntil(largeSet, k1, largeLength, smallSet[k2])
			if k1 == largeLength {
				break
			}
		}
		if smallSet[k2] < largeSet[k1] {
			k2++
			if k2 == smallLength {
				break
			}
		} else { // (set2[k2] == set1[k1])
			buffer[pos] = smallSet[k2]
			pos++
			k2++
			if k2 == smallLength {
				break
			}
			k1 = advanceUntil(largeSet, k1, largeLength, smallSet[k2])
			if k1 == largeLength {
				break
			}
		}

	}
	return pos, buffer[:pos]
}

// Find the smallest integer larger than pos such that array[pos]>= min.
// If none can be found, return length. Based on code by O. Kaser.
func advanceUntil(array []uint16, pos, length int, min uint16) int {
	lower := pos + 1

	// special handling for a possibly common sequential case
	if lower >= length || array[lower] >= min {
		return lower
	}

	spansize := 1 // could set larger  bootstrap an upper limit

	for (lower+spansize) < length && array[lower+spansize] < min {
		spansize *= 2
	}
	var upper int
	if lower+spansize < length {
		upper = lower + spansize
	} else {
		upper = length - 1
	}

	// maybe we are lucky (could be common case when the seek ahead
	// expected to be small and sequential will otherwise make us look bad)
	if array[upper] == min {
		return upper
	}

	if array[upper] < min { // means array has no item >= min
		return length
	}

	// we know that the next-smallest span was too small
	lower += (spansize / 2)

	// else begin binary search
	// invariant: array[lower]<min && array[upper]>min
	for lower+1 != upper {
		mid := (lower + upper) / 2
		if array[mid] == min {
			return mid
		} else if array[mid] < min {
			lower = mid
		} else {
			upper = mid
		}
	}
	return upper
}
