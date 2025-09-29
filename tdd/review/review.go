package review

import (
	"errors"
	"fmt"
)

type Recommendation string
type Role string

const (
	Recommend    Recommendation = "Recommend"
	NotRecommend Recommendation = "Not Recommend"
)

const (
	ReadWrite Role = "ReadWrite"
	Read      Role = "Read"
)

type Review struct {
	ID          int
	Stars       int
	Description string
	Recommend   string
}

type ReviewManager struct {
	Reviews []Review
	Role    Role
	nextID  int
}

func NewReviewManager() *ReviewManager {
	return &ReviewManager{Reviews: []Review{}, Role: ReadWrite, nextID: 1}
}

func (rm *ReviewManager) AddReview(description, recommend string, stars int) (Review, error) {
	if rm.Role == Read {
		return Review{}, fmt.Errorf("you are not allowed to review")
	}

	if stars > 6 || stars < 0 {
		return Review{}, fmt.Errorf("stars are only between 1 and 5")
	}

	if recommend != string(Recommend) && recommend != string(NotRecommend) {
		return Review{}, fmt.Errorf("recommend must be 0 (Recommend) or 1 (NotRecommend)")
	}

	review := Review{
		ID:          rm.nextID,
		Stars:       stars,
		Description: description,
		Recommend:   recommend,
	}

	rm.Reviews = append(rm.Reviews, review)
	rm.nextID++

	for _, review := range rm.Reviews {
		AddOrUpdateReview(&review)
	}
	return review, nil
}

func (rm *ReviewManager) UpdateStatus(id, stars int, description, recommend string) error {
	if rm.Role == Read {
		return fmt.Errorf("you are not allowed to review")
	}

	for i := range rm.Reviews {
		currentReview := &rm.Reviews[i]

		if currentReview.ID == id {
			changed := false

			if currentReview.Description != description && description != "" {
				currentReview.Description = description
				changed = true
			}

			if recommend != "" {
				if recommend != "Recommend" && recommend != "Not Recommend" {
					return errors.New("invalid recommendation")
				}

				if currentReview.Recommend != recommend {
					rm.Reviews[i].Recommend = recommend
					changed = true
				}
			}

			if stars != 0 {
				if stars < 1 || stars > 5 {
					return errors.New("stars must be between 1 and 5")
				}

				if currentReview.Stars != stars {
					currentReview.Stars = stars
					changed = true
				}
			}

			if !changed {
				return errors.New("nothing has changed or its the same")
			}
			for _, review := range rm.Reviews {
				AddOrUpdateReview(&review)
			}
			return nil
		}
	}
	for _, review := range rm.Reviews {
		AddOrUpdateReview(&review)
	}
	return errors.New("Review not found")
}

func (rm *ReviewManager) DeleteReview(id int) error {
	if rm.Role == Read {
		return fmt.Errorf("you are not allowed to review")
	}

	err := DeleteReviewId(id)
	if err != nil {
		return err
	}

	return nil
}

func (rm *ReviewManager) GetReviews() ([]Review, error) {
	reviews, err := LoadReview()
	if err != nil {
		return nil, err
	}
	return reviews, nil
}
