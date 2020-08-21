package roaring

import (
	"testing"
)

func TestBitmapContains(t *testing.T) {
	bc := newBitmapContainer()

	for i := 0; i < 10; i++ {
		bc.add(uint16(i))
	}

	if bc.contains(uint16(100)) {
		t.Errorf("Contains: %v, want: %v", true, false)
	}

	for i := 0; i < 10; i++ {
		contains := bc.contains(uint16(i))
		if !contains {
			t.Errorf("Contains: %v, want: %v", contains, true)
		}
	}
}

func TestAndBitmapArray(t *testing.T) {
	bt := newBitmapContainer()
	ac := newArrayContainer()

	for i := 0; i < (1 << 16); i++ {
		bt.add(uint16(i))
	}
	for i := 0; i < 5; i++ {
		ac.add(uint16(i))
	}

	answer := bt.and(ac)
	if answer.getCardinality() != 5 {
		t.Errorf("Cardinality: %d, want: 5", answer.getCardinality())
	}
	for i := 0; i < 5; i++ {
		if !answer.contains(uint16(i)) {
			t.Errorf("Container missing: %d", i)
			break
		}
	}
}
