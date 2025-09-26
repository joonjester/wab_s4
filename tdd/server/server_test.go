package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"tdd/review"
	"testing"
)

func setup() {
	rm = review.NewReviewManager()
}

func TestAddReviewHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       string
		wantStatus int
	}{
		{"Success", http.MethodPost, `{"Description":"Test","Recommend": "Not Recommend", "Stars": 5}`, http.StatusCreated},
		{"Wrong Method", http.MethodGet, ``, http.StatusMethodNotAllowed},
		{"Invalid JSON", http.MethodPost, `{bad-json`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			req := httptest.NewRequest(tt.method, "/add", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			addReviewHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestGetReviewHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{"Success", http.MethodGet, http.StatusOK},
		{"Wrong Method", http.MethodPost, http.StatusMethodNotAllowed},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			rm.AddReview("Desc", "Recommend", 5)
			req := httptest.NewRequest(tt.method, "/get", nil)
			w := httptest.NewRecorder()

			getReviewHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestUpdateReviewHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       any
		wantStatus int
	}{
		{"Success", http.MethodPut, review.Review{ID: 1, Stars: 5}, http.StatusOK},
		{"Wrong Method", http.MethodGet, nil, http.StatusMethodNotAllowed},
		{"Invalid JSON", http.MethodPut, "{bad-json", http.StatusBadRequest},
		{"Not Found", http.MethodPut, review.Review{ID: 999, Stars: 5}, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			rm.AddReview("Desc", "Not Recommend", 1)

			var bodyBytes []byte
			switch v := tt.body.(type) {
			case string:
				bodyBytes = []byte(v)
			case review.Review:
				bodyBytes, _ = json.Marshal(v)
			case nil:
				bodyBytes = nil
			}

			req := httptest.NewRequest(tt.method, "/update", bytes.NewReader(bodyBytes))
			w := httptest.NewRecorder()

			updateReviewHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}

func TestDeleteReviewHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		id         string
		wantStatus int
	}{
		{"Success", http.MethodDelete, "1", http.StatusOK},
		{"Wrong Method", http.MethodGet, "1", http.StatusMethodNotAllowed},
		{"Invalid ID", http.MethodDelete, "abc", http.StatusBadRequest},
		{"Not Found", http.MethodDelete, "999", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setup()
			rm.AddReview("Desc", "Recommend", 5)

			req := httptest.NewRequest(tt.method, "/delete?id="+tt.id, nil)
			w := httptest.NewRecorder()

			deleteReviewHandler(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("got %d, want %d", w.Code, tt.wantStatus)
			}
		})
	}
}
