package main

import (
	"fmt"
	"github.com/shrivatsas/huduku/loaders"
	"github.com/shrivatsas/huduku/search"
	"os"
)

func main() {
	docs, err := loaders.LoadDocuments("/Users/shrivatsa/Downloads/Wikipedia/enwiki-latest-abstract1.xml")
	if err != nil {
		fmt.Printf("Oh..Oh %s", err)
		os.Exit(1)
	}

	fmt.Println("Search by String Contains")
	for i, doc := range search.Str(docs, "Hello") {
		fmt.Println(i, doc.Title)
	}

	fmt.Println("Search by Regex")
	for i, doc := range search.Re(docs, "Hello") {
		fmt.Println(i, doc.Title)
	}

	fmt.Println("Search by Index")
	for i, doc := range search.Idx(docs, "Hello") {
		fmt.Println(i, doc.Title)
	}
}
