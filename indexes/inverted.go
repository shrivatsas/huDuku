package indexes

import (
	"fmt"
	"github.com/shrivatsas/huduku/loaders"
)

// Index maintains an in-memory Map of Words : [doc.ID]
type Index map[string][]int

func (idx Index) add(docs []loaders.Document) {
	for _, doc := range docs {
		for _, token := range Analyze(doc.Text) {
			ids := idx[token]
			if ids != nil && ids[len(ids)-1] == doc.ID {
				// Don't add same ID twice.
				continue
			}
			idx[token] = append(ids, doc.ID)
		}
	}
}

// CreateInverted creates the inverted index for input documents
func CreateInverted(docs []loaders.Document) *Index {
	idx := make(Index)
	idx.add(docs)
	return &idx
}

// Inverted finds all doc ids which have the given string
func Inverted(idx *Index, text string) [][]int {
	var r [][]int
	for _, token := range Analyze(text) {
		if ids, ok := (*idx)[token]; ok {
			r = append(r, ids)
		}
	}
	return r
}
