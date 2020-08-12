package persistance

import (
	"testing"
)

func TestRoot(t *testing.T) {
	root := New()
	_, ok := root.Delete("")
	if ok {
		t.Fatal("Delete on a Empty tree")
	}

	_, ok = root.Insert("", nil)
	if ok {
		t.Fatal("Inserting an Empty key")
	}

	val, ok := root.Get("")
	if !ok || val != "" {
		t.Fatal("Empty search fails")
	}
}

func TestSingleValue(t *testing.T) {
	root := New()
	a, ok := root.Insert("aa", 123)
	if a != 123 && !ok {
		t.Fatalf("Expected 123: Found %s", a)
	}

	b, ok := root.Insert("b", 456)
	if b != 456 && !ok {
		t.Fatalf("Expected 456: Found %s", b)
	}

	c, ok := root.Insert("c", 789)
	if c != 789 && !ok {
		t.Fatalf("Expected 789: Found %s", c)
	}

	aa, ok := root.Get("aa")
	if aa != 123 || !ok {
		t.Fatalf("Expected 123: Found %s", aa)
	}

	cc, ok := root.Get("c")
	if cc != 789 || !ok {
		t.Fatalf("Expected 789: Found %s", cc)
	}

	dd, ok := root.Get("d")
	if dd != nil || ok {
		t.Fatalf("Expected nil: Found %s", dd)
	}
}

func TestMultiValue(t *testing.T) {
	root := New()
	a, ok := root.Insert("aa", 123)
	if a != 123 && !ok {
		t.Fatalf("Expected 123: Found %s", a)
	}

	b, ok := root.Insert("ab", 456)
	if b != 456 && !ok {
		t.Fatalf("Expected 456: Found %s", b)
	}

	c, ok := root.Insert("bbc", 789)
	if c != 789 && !ok {
		t.Fatalf("Expected 789: Found %s", c)
	}

	ab, ok := root.Get("ab")
	if ab != 456 || !ok {
		t.Fatalf("Expected 456: Found %s", ab)
	} else {
		t.Log("found ab")
	}

	aa, ok := root.Get("aa")
	if aa != 123 || !ok {
		t.Fatalf("Expected 123: Found %s", aa)
	} else {
		t.Log("found aa")
	}

	dd, ok := root.Get("abd")
	if dd != nil || ok {
		t.Fatalf("Expected nil: Found %s", dd)
	} else {
		t.Log("not found abd")
	}
}
