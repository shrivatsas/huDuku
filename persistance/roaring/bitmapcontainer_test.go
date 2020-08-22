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

func TestAndBitmap(t *testing.T) {
	bt1 := newBitmapContainer()
	bt2 := newBitmapContainer()

	for i := 0; i < 4000; i++ {
		bt1.add(uint16(i))
		bt2.add(uint16(i))
	}

	answer := bt1.and(bt2)
	switch ac := answer.(type) {
	case *arrayContainer:
		if ac.cardinality != 4000 {
			t.Errorf("Cardinality: %d, want: 4000", ac.cardinality)
		}
		for i := 0; i < 4000; i++ {
			if ac.values[i] != uint16(i) {
				t.Errorf("AndBitmap: %d, want: %d", ac.values[i], i)
				break
			}
		}
	default:
		t.Errorf("Wrong container type: %T", ac)
	}
}

func TestOrBitmap(t *testing.T) {
	bt1 := newBitmapContainer()
	bt2 := newBitmapContainer()

	for i := 0; i < 4000; i++ {
		bt1.add(uint16(i))
		bt2.add(uint16(i))
	}

	answer := bt1.or(bt2)
	switch ac := answer.(type) {
	case *arrayContainer:
		if ac.cardinality != 4000 {
			t.Errorf("Cardinality: %d, want: 4000", ac.cardinality)
		}
		for i := 0; i < 4000; i++ {
			if ac.values[i] != uint16(i) {
				t.Errorf("AndBitmap: %d, want: %d", ac.values[i], i)
				break
			}
		}
	default:
		t.Errorf("Wrong container type: %T", ac)
	}
}
