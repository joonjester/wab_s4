package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"tdd/review"
)

var rm *review.ReviewManager

func addReviewHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var r review.Review
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	newReview, err := rm.AddReview(r.Description, r.Recommend, r.Stars)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newReview)
}

func updateReviewHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var r review.Review

	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if err := rm.UpdateStatus(r.ID, r.Stars, r.Description, r.Recommend); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func deleteReviewHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	idStr := req.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "the id is not a number", http.StatusBadRequest)
		return
	}

	if err := rm.DeleteReview(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getReviewHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rm.GetReviews())
}

func Rounter() {
	rm = review.NewReviewManager()

	http.HandleFunc("/add", addReviewHandler)
	http.HandleFunc("/get", getReviewHandler)
	http.HandleFunc("/update", updateReviewHandler)
	http.HandleFunc("/delete", deleteReviewHandler)

	http.ListenAndServe(":8080", nil)
}
