# natomanga-api

Unofficial natomanga web scraper written in Go using [gocolly](https://github.com/gocolly/colly).

## install

```
go get -u github.com/benjamindimalanta/natomanga-api
```

## usage

example:
```go
package main

import (
	"fmt"
	"log"

	nato "github.com/benjamindimalanta/natomanga-api"
)

func main() {
	searcher := nato.NewSearcher()

	mangas, err := searcher.SearchManga("Chainsaw Man")
	if err != nil {
		log.Panicln(err)
	}

	for _, manga := range *mangas {
		fmt.Println(manga.ID)
		fmt.Println(manga.Name)
		fmt.Println(manga.Author.Name)
		fmt.Println(manga.Updated)
	}

}

```

terminal output:
```
dn980422
Chainsaw Man
Tatsuki Fujimoto
Updated : Dec 14,2020 - 04:53
```

## features

1. search mangas by name
2. search specific manga by id
3. search specific manga by chapter id, returns pages of that chapter
4. search mangas by author id
5. search mangas by genre id
6. get list of the latest updated mangas (natomanga's home page)

## credits

- [gocolly](https://github.com/gocolly/colly) for web scraper
