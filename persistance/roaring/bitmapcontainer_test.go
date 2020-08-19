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
	for i := 0; i < 4096; i++ {
		ac.add(uint16(i))
	}

	answer := bt.andArray(ac)
	if answer.cardinality != 4096 {
		t.Errorf("Cardinality: %d, want: 4096", answer.cardinality)
	}
	for i := 0; i < 4096; i++ {
		if answer.values[i] != uint16(i) {
			t.Errorf("AndBitmapArray: %d, want: %d", answer.values[i], i)
			break
		}
	}
}

func TestToArrayContainer(t *testing.T) {
	bc := newBitmapContainer()

	for i := 0; i < 4096; i++ {
		bc.add(uint16(i))
	}

	ac := bc.toArrayContainer()

	if ac.cardinality != 4096 {
		t.Errorf("Cardinality: %d, want: %d", ac.cardinality, 4096)
	}
	for i := 0; i < 4096; i++ {
		if ac.values[i] != uint16(i) {
			t.Errorf("Content: %d, want: %d", ac.values[i], i)
			break
		}
	}
}
