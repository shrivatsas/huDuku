package roaring

import (
	"math/rand"
	"testing"
)

func TestToBitmapContainer(t *testing.T) {
	ac := newArrayContainer()

	rand.Seed(42)
	for i := 0; i < 5; i++ {
		ac.add(uint16(rand.Int31n(1 << 16)))
	}

	bc := arrayToBitmap(ac)
	for i := 0; i < ac.cardinality; i++ {
		if !bc.contains(ac.values[i]) {
			t.Errorf("Bitmap missing: %d", ac.values[i])
		}
	}
}

func TestToArrayContainer(t *testing.T) {
	bc := newBitmapContainer()

	for i := 0; i < 4096; i++ {
		bc.add(uint16(i))
	}

	ac := bitmapToArray(bc)
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
