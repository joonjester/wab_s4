package review

import (
	"os"
	"testing"
)

func setupTestDB(t *testing.T) func() {
	// Backup original
	originalDbFile := dbFile
	dbFile = "test_review_db.json"

	// Create empty file
	os.WriteFile(dbFile, []byte("[]"), 0644)

	// Return cleanup function
	return func() {
		os.Remove(dbFile)
		dbFile = originalDbFile
	}
}

func TestAddReview(t *testing.T) {
	rm := NewReviewManager()

	tests := map[string]struct {
		stars         int
		recommend     string
		description   string
		wantStar      int
		wantRecommend string
		wantErr       bool
	}{
		"Happy Path: Recommend": {
			stars:         5,
			recommend:     "Recommend",
			description:   "Test Test",
			wantStar:      5,
			wantRecommend: "Recommend",
			wantErr:       false,
		},
		"Happy Path: Not Recommend": {
			stars:         1,
			recommend:     "Not Recommend",
			description:   "Test Test",
			wantStar:      1,
			wantRecommend: "Not Recommend",
			wantErr:       false,
		},
		"Unhappy Path: Invalid stars": {
			stars:   6,
			wantErr: true,
		},
		"Unhappy Path: Invalid recommend": {
			recommend: "Test",
			wantErr:   true,
		},
	}

	for name, tt := range tests {
		cleanup := setupTestDB(t)
		defer cleanup()
		t.Run(name, func(t *testing.T) {
			review, err := rm.AddReview(tt.description, tt.recommend, tt.stars)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Error(err)
				return
			}

			if review.Recommend != tt.wantRecommend {
				t.Errorf("got %v, want %v", review.Recommend, tt.wantRecommend)
			}
		})
	}
}

func TestUpdateReview(t *testing.T) {
	tests := map[string]struct {
		id              int
		stars           int
		recommend       string
		description     string
		wantStars       int
		wantRecommend   string
		wantDescription string
		wantErr         bool
	}{
		"Happy Path: Stars updated": {
			id:        1,
			stars:     5,
			wantStars: 5,
			wantErr:   false,
		},
		"Happy Path: Recommend updated": {
			id:            2,
			recommend:     "Recommend",
			wantRecommend: "Recommend",
			wantErr:       false,
		},
		"Happy Path: Description updated": {
			id:              1,
			description:     "test",
			wantDescription: "test",
			wantErr:         false,
		},
		"Unhappy Path: Review not found": {
			id:      5,
			wantErr: true,
		},
		"Unhappy Path: Invalid stars": {
			id:      1,
			stars:   6,
			wantErr: true,
		},
		"Unhappy Path: Invalid recommend": {
			id:        1,
			recommend: "Test",
			wantErr:   true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup fresh database for each test
			cleanup := setupTestDB(t)
			defer cleanup()

			// Create new ReviewManager and add test data
			rm := NewReviewManager()
			rm.AddReview("Beschreibung", "Recommend", 0)
			rm.AddReview("Beschreibung2", "Not Recommend", 1)
			rm.AddReview("Beschreibung3", "Recommend", 0)

			// Run the update
			err := rm.UpdateStatus(tt.id, tt.stars, tt.description, tt.recommend)

			if tt.wantErr {
				if err == nil {
					t.Error("wanted error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Get reviews and check results
			reviews, err := rm.GetReviews()
			if err != nil {
				t.Fatalf("GetReviews failed: %v", err)
			}

			if len(reviews) < tt.id {
				t.Fatalf("Not enough reviews: got %d, need at least %d", len(reviews), tt.id)
			}

			// Check the updated review
			updatedReview := reviews[tt.id-1]

			if tt.wantDescription != "" {
				if updatedReview.Description != tt.wantDescription {
					t.Errorf("Description: got %v, want %v", updatedReview.Description, tt.wantDescription)
				}
			}
			if tt.wantStars != 0 {
				if updatedReview.Stars != tt.wantStars {
					t.Errorf("Stars: got %v, want %v", updatedReview.Stars, tt.wantStars)
				}
			}
			if tt.wantRecommend != "" {
				if updatedReview.Recommend != tt.wantRecommend {
					t.Errorf("Recommend: got %v, want %v", updatedReview.Recommend, tt.wantRecommend)
				}
			}
		})
	}
}

func TestDeleteReview(t *testing.T) {
	tests := map[string]struct {
		id      int
		wantErr bool
	}{
		"Happy Path: could find the Review": {
			id:      1,
			wantErr: false,
		},
		"Unhappy Path: couldn't find the Review": {
			id:      2,
			wantErr: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			// Setup fresh database for each test
			cleanup := setupTestDB(t)
			defer cleanup()

			rm := NewReviewManager()

			// Add only ONE review with ID 1
			_, _ = rm.AddReview("Beschreibung", "Recommend", 5)

			// Try to delete
			err := rm.DeleteReview(tt.id)

			if tt.wantErr {
				if err == nil {
					t.Error("wanted error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			reviews, err := rm.GetReviews()
			if err != nil {
				t.Fatalf("GetReviews() failed: %v", err)
			}

			if len(reviews) != 0 {
				t.Errorf("expected no reviews, got %d: %v", len(reviews), reviews)
			}
		})
	}
}
