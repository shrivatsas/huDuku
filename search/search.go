package search

import (
	"fmt"
	"github.com/shrivatsas/huduku/indexes"
	"github.com/shrivatsas/huduku/loaders"
	"regexp"
	"strings"
)

var idx indexes.Index

// Str for a string among the provided documents
func Str(docs []loaders.Document, term string) []loaders.Document {
	var r []loaders.Document
	for _, doc := range docs {
		if strings.Contains(doc.Text, term) {
			r = append(r, doc)
		}
	}
	return r
}

// Re for a string among the provided documents
func Re(docs []loaders.Document, term string) []loaders.Document {
	re := regexp.MustCompile(`(?i)\b` + term + `\b`)
	var r []loaders.Document
	for _, doc := range docs {
		if re.MatchString(doc.Text) {
			r = append(r, doc)
		}
	}
	return r
}

// Idx finds a string in provided Index, and returns doc.IDs
func Idx(docs []loaders.Document, term string) []loaders.Document {
	if idx == nil {
		fmt.Println("Creating Index")
		idx = *indexes.CreateInverted(docs)
	}
	var r []loaders.Document
	docIds := indexes.Inverted(&idx, term)
	for _, id := range docIds {
		r = append(r, docs[id])
	}
	return r
}
