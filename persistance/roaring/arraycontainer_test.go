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

func TestNewArrayContainerRunOfOnes(t *testing.T) {
	ac := newArrayContainerRunOfOnes(1, 4)
	if ac.cardinality != 4 {
		t.Errorf("Cardinality = %d, want %d", ac.cardinality, 4)
	}
	if (ac.values[0] != 1) ||
		(ac.values[1] != 2) ||
		(ac.values[2] != 3) ||
		(ac.values[3] != 4) {
		t.Errorf("Content error = %v", ac.values)
	}
}

func TestArrayContainerAdd(t *testing.T) {
	ac := newArrayContainer()
	ac.add(uint16(0))
	ac.add(uint16(2))

	if ac.values[0] != uint16(0) &&
		ac.values[1] != uint16(2) {
		t.Errorf("Wrong add: %d, %d, want: %d, %d", ac.values[0], ac.values[1],
			0, 2)
	}

	if ac.cardinality != 2 {
		t.Errorf("Cardinality: %d, want %d", ac.cardinality, 2)
	}

	ac.add(uint16(1))
	if ac.values[0] != uint16(0) &&
		ac.values[1] != uint16(1) &&
		ac.values[2] != uint16(2) {
		t.Errorf("Wrong add: %d, %d %d, want: %d, %d %d",
			ac.values[0],
			ac.values[1],
			ac.values[2],
			0, 1, 2)
	}

	if ac.cardinality != 3 {
		t.Errorf("Cardinality: %d, want %d", ac.cardinality, 3)
	}
}

func TestArrayContains(t *testing.T) {
	ac := newArrayContainerRunOfOnes(0, 9)

	for i := 10; i < 20; i++ {
		if ac.contains(uint16(i)) {
			t.Errorf("Array constains %d, want: false)", i)
		}
	}

	for i := 0; i < 9; i++ {
		if !ac.contains(uint16(i)) {
			t.Errorf("Array not constains %d, want: true)", i)
		}
	}

	ac = newArrayContainerRunOfOnes(1, 5)

	if ac.contains(uint16(0)) {
		t.Errorf("Array constains zero.")
	}
}

func TestOrArray(t *testing.T) {
	ac1 := newArrayContainer()
	ac2 := newArrayContainer()
	count := 100

	for i := 0; i < count; i += 2 {
		ac1.add(uint16(i))
		ac2.add(uint16(i + 1))
	}

	result := ac1.orArray(ac2)
	if result.getCardinality() != count {
		t.Errorf("Cardinality: %d, want: %d", result.getCardinality(), count)
	}
	for k, v := range result.(*arrayContainer).values {
		if v != uint16(k) {
			t.Errorf("orArray: %d, want: %d", v, k)
		}
	}
}
