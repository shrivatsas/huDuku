package main

import (
	"fmt"
	"github.com/shrivatsas/huduku/loaders"
	"github.com/shrivatsas/huduku/search"
	"os"
	"strings"
)

func main() {
	docs, err := loaders.LoadDocuments("/Users/shrivatsa/Downloads/Wikipedia/enwiki-latest-abstract1.xml")
	if err != nil {
		fmt.Printf("Oh..Oh %s", err)
		os.Exit(1)
	}

	argsWithoutProg := strings.Join(os.Args[1:], " ")
	fmt.Printf("Searching for %s\n", argsWithoutProg)
	fmt.Println("Search by String Contains")
	for i, doc := range search.Str(docs, argsWithoutProg) {
		fmt.Println(i, doc.Title)
	}

	fmt.Println("Search by Regex")
	for i, doc := range search.Re(docs, argsWithoutProg) {
		fmt.Println(i, doc.Title)
	}

	fmt.Println("Search by Index")
	for i, doc := range search.Idx(docs, argsWithoutProg) {
		fmt.Println(i, doc.Title)
	}
}
