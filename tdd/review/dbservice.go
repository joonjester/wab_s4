package review

import (
	"encoding/json"
	"fmt"
	"os"
)

var dbFile = "review_db.json"

func AddOrUpdateReview(newReview *Review) {
	allReview, err := LoadReview()
	if err != nil {
		fmt.Println("could not load review:", err)
	}

	updated := false
	for index, review := range allReview {
		if review.ID == newReview.ID {
			allReview[index] = *newReview
			updated = true
			break
		}
	}

	if !updated {
		allReview = append(allReview, *newReview)
	}

	data, err := json.MarshalIndent(allReview, "", "  ")
	if err != nil {
		fmt.Println("could not save review:", err)
		return
	}

	if err := os.WriteFile(dbFile, data, 0644); err != nil {
		fmt.Println("could not write file:", err)
	}
}

func LoadReview() ([]Review, error) {
	data, err := os.ReadFile(dbFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Review{}, nil
		}
		return nil, err
	}
	var users []Review
	err = json.Unmarshal(data, &users)
	return users, err
}

func DeleteReviewId(id int) error {
	allReviews, err := LoadReview()
	if err != nil {
		return fmt.Errorf("could not load reviews: %w", err)
	}

	indexToDelete := -1
	for i, review := range allReviews {
		if review.ID == id {
			indexToDelete = i
			break
		}
	}

	if indexToDelete == -1 {
		return fmt.Errorf("review with ID %d not found", id)
	}

	allReviews = append(allReviews[:indexToDelete], allReviews[indexToDelete+1:]...)

	data, err := json.MarshalIndent(allReviews, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal reviews: %w", err)
	}

	if err := os.WriteFile(dbFile, data, 0644); err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil

}
