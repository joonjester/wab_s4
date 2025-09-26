package review

import (
	"testing"
)

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
	rm := NewReviewManager()

	rm.AddReview("Beschreibung", "Recommend", 0)
	rm.AddReview("Beschreibung2", "Not Recommend", 1)
	rm.AddReview("Beschreibung3", "Recommend", 0)

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

			if tt.wantDescription != "" {
				got := rm.GetReviews()[tt.id-1].Description
				if got != tt.wantDescription {
					t.Errorf("got %v, want %v", got, tt.wantDescription)
					return
				}
			}

			if tt.wantStars != 0 {
				got := rm.GetReviews()[tt.id-1].Stars
				if got != tt.wantStars {
					t.Errorf("got %v, want %v", got, tt.stars)
					return
				}
			}

			if tt.wantRecommend != "" {
				got := rm.GetReviews()[tt.id-1].Recommend
				if got != tt.wantRecommend {
					t.Errorf("got %v, want %v", got, tt.wantRecommend)
					return
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
			rm := NewReviewManager()
			_, _ = rm.AddReview("Beschreibung", "Recommend", 5)

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

			if len(rm.GetReviews()) != 0 {
				t.Errorf("expected no Reviews, got %d", len(rm.GetReviews()))
				return
			}
		})
	}
}
