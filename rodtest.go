package main

import (
    "fmt"
    "github.com/go-rod/rod"
)

func main() {
    page := rod.New().MustConnect().MustPage("https://www.natomanga.com")
    title := page.MustEval("() => document.title")
    fmt.Println("Page title:", title)
}
