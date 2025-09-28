package books

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Book struct {
	Title         string   `json:"title"`
	PublishedYear int      `json:"first_publish_year"`
	AuthorName    []string `json:"author_name"`
	Language      []string `json:"language"`
}

type SearchResult struct {
	NumFound int    `json:"numFound"`
	Start    int    `json:"start"`
	Books    []Book `json:"docs"`
}

var bookAPI = "https://openlibrary.org/search.json"

func SearchBooks(title string) (*Book, error) {
	url := fmt.Sprintf("%s?q=%v", bookAPI, title)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("wrong title: %s", title)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SearchResult
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	for _, books := range result.Books {
		return &books, nil
	}

	return nil, fmt.Errorf("No book found called: %s", title)
}
