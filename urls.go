package natomangaapi

const (
	natomangaURL              string = "natomanga.com"
	natomangaURLWithHTTPS            = "https://" + natomangaURL
	readnatomangaURL                 = "chap" + natomangaURL
	readnatomangaURLWihtHTTPS        = "https://" + readnatomangaURL
	searchMangaURL                   = natomangaURLWithHTTPS + "/search/story/"
	specificMangaURL                 = readnatomangaURLWihtHTTPS + "/manga-"
	searchMangaByAuthorURL           = natomangaURLWithHTTPS + "/author/story/"
	searchMangaByGenreURL            = natomangaURLWithHTTPS + "/genre-"
)
