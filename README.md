## Overview
### huDuku (ಹುಡುಕು)
A toy project to learn golang as well as some text search techniques and optimizations.

Further ideas from Artem's blog
1. Extend boolean queries to support OR and NOT.
2. Store the index on disk:
    Rebuilding the index on every application restart may take a while.
    Large indexes may not fit in memory.
3. Experiment with memory and CPU-efficient data formats for storing sets of document IDs.
4. Support indexing multiple document fields.
5. Sort results by relevance, using a score.

Others
1. Language support (Indian languages)

The roadmap is managed via Github Projects.

## Running a search
```
go run main.go hello world
```
The above command tries to run a search for the words "hello world"

## References
https://artem.krylysov.com/blog/2020/07/28/lets-build-a-full-text-search-engine/  
https://github.com/akrylysov/simplefts  
https://github.com/fzandona/goroar  
