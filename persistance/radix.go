package persistance

import (
	"fmt"
)

type node struct {
	// the label / edge connecting current node with its parent
	label []byte

	// the child nodes
	children []*node

	// the leaf value, if current node is a leaf
	value interface{}
}

// Tree is the implementation of a Radix tree
type Tree struct {
	root *node
	size int
}

// New creates a new Radix tree with an empty root
func New() *Tree {
	return &Tree{root: &node{}}
}

// Delete deletes the provided key, and returns the previous value
func (t *Tree) Delete(key string) (interface{}, bool) {
	if len(key) == 0 {
		return nil, false
	}
	return nil, true
}

// Insert inserts a new key and val into the radix tree
func (t *Tree) Insert(key string, val interface{}) (interface{}, bool) {
	fmt.Printf("\nInserting %s\n", key)
	if len(key) == 0 {
		return nil, false
	}

	parent := t.root
	var i int
	ok := false
	search := []byte(key)
	for {
		if len(search) > 0 {
			parent, i, ok = parent.longestPrefix(search)
			fmt.Printf("Outcome %s %d %t\n", parent.label, i, ok)
			if !ok && i == 0 {
				n := &node{label: search, value: val}
				parent.children = append(parent.children, n)
				fmt.Printf("Added single child %s\n", n.label)
				break
			}

			if !ok {
				fmt.Println("Split children ")
				if len(parent.label) > i {
					n1 := &node{label: parent.label[(i - 1):], value: parent.value}
					parent.label = parent.label[:i]
					parent.value = nil
					parent.children = append(parent.children, n1)
					fmt.Printf("> C1 %s\n", n1.label)
				}
				n2 := &node{label: search[(i):], value: val}
				parent.children = append(parent.children, n2)
				fmt.Printf("> C2 %s\n", n2.label)
				fmt.Printf("> Parent %s\n", parent.label)
				break
			}
			search = search[i:]
		}
	}
	return key, true
}

// Get fetched a value for the provided key
func (t *Tree) Get(key string) (interface{}, bool) {
	fmt.Printf("\nGetting %s\n", key)
	if len(key) == 0 {
		return "", true
	}

	var val interface{}
	parent := t.root
	var i int
	ok := false
	search := []byte(key)
	for {
		parent, i, ok = parent.longestPrefix(search)
		fmt.Printf("Searched: %s, Outcome %s %d %t\n", search, parent.label, i, ok)
		if !ok && i == 0 {
			val = nil
			break
		}

		search = search[i:]
		if len(search) == 0 {
			val = parent.value
			ok = true
			break
		}
	}

	fmt.Printf("Got %s %t\n", val, ok)
	return val, ok
}

func (n *node) longestPrefix(s []byte) (*node, int, bool) {
	fmt.Printf("Current: Node %s, Key %s\n", n.label, s)
	child := n
	var matched int
	ok := false
	// search := s[len(n.label):]
	for _, nd := range n.children {
		l := nd.childPrefix(s)
		if matched < l {
			matched = l
			child = nd
			ok = true
			break
		}
	}

	return child, matched, ok && len(child.children) > 0
}

func (n *node) childPrefix(s []byte) int {
	fmt.Printf("> Child: Node %s, Key %s\n", n.label, s)
	maxLen := len(n.label)
	if maxLen > len(s) {
		maxLen = len(s)
	}

	var i int
	for i < maxLen {
		if n.label[i] != s[i] {
			break
		}
		i++
	}
	fmt.Printf("%s -- %s == %d\n", n.label, s, i)
	return i
}
