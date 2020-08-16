package roaring

import (
	"testing"
)

func TestContains(t *testing.T) {
	bc := newBitmapContainer()

	for i := 0; i < 10; i++ {
		bc.add(uint16(i))
	}

	// if bc.contains(uint16(100)) {
	// 	t.Errorf("Contains: %v, want: %v", true, false)
	// }

	// for i := 0; i < 10; i++ {
	// 	contains := bc.contains(uint16(i))
	// 	if !contains {
	// 		t.Errorf("Contains: %v, want: %v", contains, true)
	// 	}
	// }
}
