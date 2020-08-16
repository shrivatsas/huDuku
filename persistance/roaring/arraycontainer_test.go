package roaring

import (
	"testing"
)

func TestNewArrayContainer(t *testing.T) {
	ac := newArrayContainer()
	if ac.cardinality != 0 {
		t.Errorf("Cardinality = %d, want %d", ac.cardinality, 0)
	}
	if cap(ac.values) != arrayContainerInitSize {
		t.Errorf("Content = %d, want %d", cap(ac.values), arrayContainerInitSize)
	}
}

func TestNewArrayContainerWithCapacity(t *testing.T) {
	capacity := 42
	ac := newArrayContainerWithCapacity(capacity)
	if ac.cardinality != 0 {
		t.Errorf("Cardinality = %d, want %d", ac.cardinality, 0)
	}
	if cap(ac.values) != capacity {
		t.Errorf("Content = %d, want %d", cap(ac.values), capacity)
	}
}
