package books

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSearchBook(t *testing.T) {
	tests := map[string]struct {
		title      string
		wantErr    bool
		wantResult Book
	}{
		"Happy Path": {
			title:   "harry+potter",
			wantErr: false,
			wantResult: Book{
				Title:         "Harry Potter and the Philosopher's StoneAuthors",
				PublishedYear: 1997,
				AuthorName:    []string{"J. K. Rowling"},
				Language:      []string{"gre, eng"},
			},
		},
		"Unhappy Path: Book not found": {
			title:   "/////",
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.wantErr {
					response := map[string]any{
						"numFound": 0,
						"start":    0,
						"docs":     []Book{},
					}
					_ = json.NewEncoder(w).Encode(response)
					return
				}
				response := map[string]any{
					"numFound": 1,
					"start":    0,
					"docs": []map[string]any{
						{
							"title":              tt.wantResult.Title,
							"first_publish_year": tt.wantResult.PublishedYear,
							"author_name":        tt.wantResult.AuthorName,
							"language":           tt.wantResult.Language,
						},
					},
				}
				_ = json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			oldApi := bookAPI
			bookAPI = server.URL
			defer func() { bookAPI = oldApi }()

			got, err := SearchBooks(tt.title)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexcepted error: %v", err)
				return
			}

			if got == nil {
				t.Errorf("wanted book, got nil")
				return
			}

			err = compareBooks(got, tt.wantResult)
			if err != nil {
				t.Error(err)
				return
			}
		})
	}
}

func compareBooks(receivedBook *Book, wantBook Book) error {
	if receivedBook.Title != wantBook.Title {
		return fmt.Errorf("got %v, want %v", receivedBook.Title, wantBook.Title)
	}
	if receivedBook.PublishedYear != wantBook.PublishedYear {
		return fmt.Errorf("got %v, want %v", receivedBook.PublishedYear, wantBook.PublishedYear)
	}

	for i, author := range wantBook.AuthorName {
		receivedAuthor := receivedBook.AuthorName[i]
		if receivedAuthor != author {
			return fmt.Errorf("got %v, want %v", receivedAuthor, author)
		}
	}

	for i, language := range wantBook.Language {
		receivedLanguage := receivedBook.Language[i]
		if receivedLanguage != language {
			return fmt.Errorf("got %v, want %v", receivedLanguage, language)
		}
	}

	return nil
}
