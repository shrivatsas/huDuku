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

func (ac *arrayContainer) toBitmapContainer() *bitmapContainer {
	bc := newBitmapContainer()
	bc.loadData(ac)
	return bc
}
