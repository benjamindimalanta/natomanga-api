package natomangaapi

import (
	"errors"
	"sync"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
)

var (

	// throws when no results are found
	ErrPageNotFound = errors.New("this page does not exist or has been deleted")
	c               *colly.Collector
	wg              sync.WaitGroup
)

// provides methods for fetching data from https://natomanga.com
type Searcher struct {
	MethodsDescription map[string]string
}

// any struct types with ID implement Searchable
type Searchable interface {
	getID(string)
}

func initCrawler() {
func initCrawler() {
    c = colly.NewCollector(
        colly.AllowedDomains(
            "natomanga.com",
            "www.natomanga.com",
            "readnatomanga.com",
            "www.readnatomanga.com",
        ),
        colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) " +
            "AppleWebKit/537.36 (KHTML, like Gecko) " +
            "Chrome/114.0.0.0 Safari/537.36"),
        colly.MaxDepth(2),
        colly.Async(true),
    )

    // ⬇️ ADD THIS RIGHT HERE
    c.OnRequest(func(r *colly.Request) {
        r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) "+
            "AppleWebKit/537.36 (KHTML, like Gecko) "+
            "Chrome/114.0.0.0 Safari/537.36")
        r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
        r.Headers.Set("Accept-Language", "en-US,en;q=0.5")
        r.Headers.Set("Referer", "https://www.natomanga.com/")
    })

    // Optional: you can still add logging or retries here
    // extensions.RandomUserAgent(c) ← not needed anymore
}

)



	extensions.RandomUserAgent(c)
}

func cloneCrawler() *colly.Collector {
	c2 := c.Clone()
	// extensions.RandomUserAgent(c2)

	return c2
}

func deleteCrawler() {
	c = nil
}

// returns a new Searcher
func NewSearcher() Searcher {
	methodDescription := map[string]string{
		"SearchManga":              "receives name of a manga user wants to search for and returns a list of mangas that match the name",
		"PickManga":                "receives the id of the specific manga then returns that manga if found",
		"ReadMangaChapter":         "receives the manga id and chapter id then returns pages of that specific chapter",
		"PickAuthor":               "receives the id of the author then returns a list of mangas by him/her",
		"PickGenre":                "receives genre id then returns a list of mangas with that genre",
		"SearchLatestUpdatedManga": "returns list of latest updated mangas from the first page of https://natomanga.com",
		"IsSearchable":             "returns whether the struct type implements Searchable interface",
	}

	return Searcher{
		MethodsDescription: methodDescription,
	}
}

// receives name of a manga user wants to search for
// and returns a list of mangas that match the name.
// Each manga will have ID, Name, Author and Updated fields
//
// https://natomanga.com/search/story/<name>
func (s *Searcher) SearchManga(name string) (*[]Manga, error) {
	initCrawler()
	defer deleteCrawler()

	tmp := getMangaList(changeSpaceToUnderscore(name))

	if len(tmp) == 0 {
		return nil, ErrPageNotFound
	}

	mgs := []Manga{}

	for _, mg := range tmp {
		wg.Add(1)

		go func(m Manga, c2 *colly.Collector) {
			defer wg.Done()

			createAuthor(&m, c2)
			mgs = append(mgs, m)

		}(mg, cloneCrawler())
	}

	wg.Wait()

	return &mgs, nil
}

// receives the id of the specific manga
// then returns that manga if found.
// Every fields in Manga struct will be setted
//
// https://readnatomanga.com/manga-<id>
func (s *Searcher) PickManga(id string) (*Manga, error) {
	initCrawler()
	defer deleteCrawler()

	m := Manga{
		ID: id,
	}
	tmp := m

	m.getMangaByID()

	if m.compareManga(&tmp) {
		return nil, ErrPageNotFound
	}

	return &m, nil
}

// receives the manga id and chapter id
// then returns pages of that specific chapter.
// Each page will have ID and ImageURL setted
//
// https://readnatomanga.com/manga-<mangaId>/chapter-<chapterId>
func (s *Searcher) ReadMangaChapter(mangaId, chapterId string) (*[]Page, error) {
	initCrawler()
	defer deleteCrawler()

	ch := Chapter{
		ID:      chapterId,
		MangaID: mangaId,
	}

	ch.getChapterByID()

	if len(ch.Pages) == 0 {
		return nil, ErrPageNotFound
	}

	return &ch.Pages, nil
}

// receives the id of the author then returns a list of mangas by him/her.
// Each manga will have ID, Name, Author and Updated fields
//
// https://natomanga.com/author/story/<authorId>
func (s *Searcher) PickAuthor(authorId string) (*[]Manga, error) {
	initCrawler()
	defer deleteCrawler()

	a := Author{
		ID: authorId,
	}
	a.getMangaListByAuthorID()

	if len(a.Mangas) == 0 {
		return nil, ErrPageNotFound
	}

	for i, mg := range a.Mangas {
		wg.Add(1)
		go func(m Manga, index int, c2 *colly.Collector) {
			defer wg.Done()

			createAuthor(&m, c2)
			a.Mangas[index] = m
		}(mg, i, cloneCrawler())
	}

	wg.Wait()

	return &a.Mangas, nil
}

// receives genre id then returns a list of mangas with that genre.
// Each manga will have ID, Name, Author(Name only), Updated,
//
//	Views and Description fields
//
// https://natomanga.com/genre-<genreId>
func (s *Searcher) PickGenre(genreId string) (*[]Manga, error) {
	initCrawler()
	defer deleteCrawler()

	g := Genre{
		ID: genreId,
	}

	g.getMangaListByGenreID()

	if len(g.Mangas) == 0 {
		return nil, ErrPageNotFound
	}

	return &g.Mangas, nil
}

// returns list of latest updated mangas
// from the first page of https://natomanga.com.
// Each manga will have ID, Name and Author(Name only) fields
func (s *Searcher) SearchLatestUpdatedManga() (*[]Manga, error) {
	initCrawler()
	defer deleteCrawler()

	tmp := getLatestUpdatedManga()

	if len(tmp) == 0 {
		return nil, ErrPageNotFound
	}

	return &tmp, nil
}

// returns whether the struct type implements Searchable interface
func (s *Searcher) IsSearchable(any interface{}) bool {
	switch any.(type) {
	case Searchable:
		return true
	}
	return false
}
