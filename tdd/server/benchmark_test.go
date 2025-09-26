package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"tdd/review"
	"testing"
)

func BenchmarkAddReviewHandler(b *testing.B) {
	rm = review.NewReviewManager()

	payload := review.Review{Description: "desc", Recommend: "Recommend", Stars: 5}
	data, _ := json.Marshal(payload)

	b.ResetTimer() // reset before measurement
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPost, "/add", bytes.NewReader(data))
		w := httptest.NewRecorder()

		addReviewHandler(w, req)

		if w.Code != http.StatusCreated {
			b.Fatalf("unexpected status: got %d", w.Code)
		}
	}
}

func BenchmarkGetReviewHandler(b *testing.B) {
	rm = review.NewReviewManager()
	// Fülle ReviewManager einmal vorab
	rm.AddReview("Desc", "Not Recommend", 1)

	req := httptest.NewRequest(http.MethodGet, "/get", nil)

	b.ResetTimer()
	for b.Loop() {
		w := httptest.NewRecorder()
		getReviewHandler(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("unexpected status: got %d", w.Code)
		}
	}
}

func BenchmarkUpdateReviewHandler(b *testing.B) {
	rm = review.NewReviewManager()
	r, _ := rm.AddReview("Desc", "Recommend", 5)

	payload := review.Review{ID: r.ID, Recommend: r.Recommend, Stars: r.Stars}
	data, _ := json.Marshal(payload)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodPut, "/update", bytes.NewReader(data))
		w := httptest.NewRecorder()

		updateReviewHandler(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("unexpected status: got %d", w.Code)
		}
	}
}

func BenchmarkDeleteReviewHandler(b *testing.B) {
	rm = review.NewReviewManager()
	r, _ := rm.AddReview("Desc", "Recommend", 5)

	url := "/delete?id=" + strconv.Itoa(r.ID)

	b.ResetTimer()
	for b.Loop() {
		req := httptest.NewRequest(http.MethodDelete, url, nil)
		w := httptest.NewRecorder()

		deleteReviewHandler(w, req)

		if w.Code != http.StatusOK {
			b.Fatalf("unexpected status: got %d", w.Code)
		}
	}
}

func BenchmarkFullApplication(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("/add", addReviewHandler)
	mux.HandleFunc("/get", getReviewHandler)
	mux.HandleFunc("/delete", deleteReviewHandler)
	mux.HandleFunc("/update", updateReviewHandler)

	server := httptest.NewServer(mux)
	defer server.Close()

	reviewed := review.Review{Description: "Test the test", Recommend: "Recommend", Stars: 5}
	reviewedNotRecomment := review.Review{Description: "Test the test", Recommend: "Not Recommend", Stars: 1}
	reviewedData, _ := json.Marshal(reviewed)
	updatedReview, _ := json.Marshal(reviewedNotRecomment)

	b.ResetTimer()
	for b.Loop() {
		resp, _ := http.Post(server.URL+"/add", "application/json", bytes.NewBuffer(reviewedData))
		resp.Body.Close()

		req, _ := http.NewRequest(http.MethodPut, server.URL+"/update", bytes.NewBuffer(updatedReview))
		req.Header.Set("Content-Type", "application/json")
		resp, _ = http.DefaultClient.Do(req)
		resp.Body.Close()

		resp, _ = http.Get(server.URL + "/get")
		resp.Body.Close()

		req, _ = http.NewRequest(http.MethodDelete, server.URL+"/delete?id=1", nil)
		resp, _ = http.DefaultClient.Do(req)
		resp.Body.Close()

	}
}
