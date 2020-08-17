package roaring

type container interface {
	add(x uint16) container
	// and(x container) container
	// or(x container) container
	// andNot(x container) container
	// xor(x container) container

	// clone() container
	contains(x uint16) bool

	// getCardinality() int
	// sizeInBytes() int
}
