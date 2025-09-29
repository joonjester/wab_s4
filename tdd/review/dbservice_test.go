package review

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestLoadReview(t *testing.T) {
	tests := []struct {
		Description    string
		setupFile      func() error
		expectedResult []Review
		expectError    bool
		cleanup        func()
	}{
		{
			Description: "load empty file returns empty slice",
			setupFile: func() error {
				return os.WriteFile("test_review_db.json", []byte("[]"), 0644)
			},
			expectedResult: []Review{},
			expectError:    false,
			cleanup: func() {
				os.Remove("test_review_db.json")
			},
		},
		{
			Description: "load file with reviews",
			setupFile: func() error {
				reviews := []Review{
					{ID: 1, Description: "Review 1"},
					{ID: 2, Description: "Review 2"},
				}
				data, _ := json.MarshalIndent(reviews, "", "  ")
				return os.WriteFile("test_review_db.json", data, 0644)
			},
			expectedResult: []Review{
				{ID: 1, Description: "Review 1"},
				{ID: 2, Description: "Review 2"},
			},
			expectError: false,
			cleanup: func() {
				os.Remove("test_review_db.json")
			},
		},
		{
			Description: "file does not exist returns empty slice",
			setupFile: func() error {
				// Don't create file
				return nil
			},
			expectedResult: []Review{},
			expectError:    false,
			cleanup:        func() {}, // Nothing to cleanup
		},
		{
			Description: "invalid JSON returns error",
			setupFile: func() error {
				return os.WriteFile("test_review_db.json", []byte("invalid json"), 0644)
			},
			expectedResult: nil,
			expectError:    true,
			cleanup: func() {
				os.Remove("test_review_db.json")
			},
		},
	}

	// Backup original dbFile
	originalDbFile := dbFile
	defer func() {
		dbFile = originalDbFile
	}()

	for _, tt := range tests {
		t.Run(tt.Description, func(t *testing.T) {
			// Set test database file
			dbFile = "test_review_db.json"

			// Setup
			if err := tt.setupFile(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			defer tt.cleanup()

			// Execute
			result, err := LoadReview()

			// Verify error expectation
			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			// Verify result
			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("Expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestAddOrUpdateReview(t *testing.T) {
	tests := []struct {
		Description    string
		setupFile      func() error
		inputReview    *Review
		expectedResult []Review
		cleanup        func()
	}{
		{
			Description: "add new review to empty file",
			setupFile: func() error {
				return os.WriteFile("test_review_db.json", []byte("[]"), 0644)
			},
			inputReview: &Review{ID: 1, Description: "New Review"},
			expectedResult: []Review{
				{ID: 1, Description: "New Review"},
			},
			cleanup: func() {
				os.Remove("test_review_db.json")
			},
		},
		{
			Description: "add new review to existing reviews",
			setupFile: func() error {
				reviews := []Review{
					{ID: 1, Description: "Existing Review"},
				}
				data, _ := json.MarshalIndent(reviews, "", "  ")
				return os.WriteFile("test_review_db.json", data, 0644)
			},
			inputReview: &Review{ID: 2, Description: "New Review"},
			expectedResult: []Review{
				{ID: 1, Description: "Existing Review"},
				{ID: 2, Description: "New Review"},
			},
			cleanup: func() {
				os.Remove("test_review_db.json")
			},
		},
		{
			Description: "update existing review",
			setupFile: func() error {
				reviews := []Review{
					{ID: 1, Description: "Old Name"},
					{ID: 2, Description: "Another Review"},
				}
				data, _ := json.MarshalIndent(reviews, "", "  ")
				return os.WriteFile("test_review_db.json", data, 0644)
			},
			inputReview: &Review{ID: 1, Description: "Updated Name"},
			expectedResult: []Review{
				{ID: 1, Description: "Updated Name"},
				{ID: 2, Description: "Another Review"},
			},
			cleanup: func() {
				os.Remove("test_review_db.json")
			},
		},
		{
			Description: "add review when file doesn't exist",
			setupFile: func() error {
				// Don't create file
				return nil
			},
			inputReview: &Review{ID: 1, Description: "First Review"},
			expectedResult: []Review{
				{ID: 1, Description: "First Review"},
			},
			cleanup: func() {
				os.Remove("test_review_db.json")
			},
		},
	}

	// Backup original dbFile
	originalDbFile := dbFile
	defer func() {
		dbFile = originalDbFile
	}()

	for _, tt := range tests {
		t.Run(tt.Description, func(t *testing.T) {
			// Set test database file
			dbFile = "test_review_db.json"

			// Setup
			if err := tt.setupFile(); err != nil {
				t.Fatalf("Setup failed: %v", err)
			}
			defer tt.cleanup()

			// Execute
			AddOrUpdateReview(tt.inputReview)

			// Verify result by loading the file
			result, err := LoadReview()
			if err != nil {
				t.Fatalf("Failed to load reviews after update: %v", err)
			}

			if !reflect.DeepEqual(result, tt.expectedResult) {
				t.Errorf("Expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

// Benchmark tests
func BenchmarkLoadReview(b *testing.B) {
	// Setup test file with some data
	reviews := make([]Review, 1000)
	for i := 0; i < 1000; i++ {
		reviews[i] = Review{ID: i, Description: fmt.Sprintf("Review %d", i)}
	}
	data, _ := json.MarshalIndent(reviews, "", "  ")
	os.WriteFile("bench_review_db.json", data, 0644)
	defer os.Remove("bench_review_db.json")

	// Backup and set test file
	originalDbFile := dbFile
	dbFile = "bench_review_db.json"
	defer func() {
		dbFile = originalDbFile
	}()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		LoadReview()
	}
}

func BenchmarkAddOrUpdateReview(b *testing.B) {
	// Setup test file
	os.WriteFile("bench_review_db.json", []byte("[]"), 0644)
	defer os.Remove("bench_review_db.json")

	// Backup and set test file
	originalDbFile := dbFile
	dbFile = "bench_review_db.json"
	defer func() {
		dbFile = originalDbFile
	}()

	testReview := &Review{ID: 1, Description: "Benchmark Review"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AddOrUpdateReview(testReview)
	}
}
