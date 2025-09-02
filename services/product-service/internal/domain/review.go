package domain

import (
	"time"
)

// ProductReview represents a product review
type ProductReview struct {
	ID           int64     `json:"id" db:"id"`
	ProductID    int64     `json:"product_id" db:"product_id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	OrderID      *int64    `json:"order_id" db:"order_id"` // optional, for verified purchases
	Rating       int       `json:"rating" db:"rating"`     // 1-5 stars
	Title        string    `json:"title" db:"title"`
	Review       string    `json:"review" db:"review"`
	IsVerified   bool      `json:"is_verified" db:"is_verified"` // verified purchase
	IsApproved   bool      `json:"is_approved" db:"is_approved"`
	IsHelpful    int       `json:"is_helpful" db:"is_helpful"`         // helpful votes count
	IsNotHelpful int       `json:"is_not_helpful" db:"is_not_helpful"` // not helpful votes count
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// ReviewImage represents images attached to reviews
type ReviewImage struct {
	ID       int64  `json:"id" db:"id"`
	ReviewID int64  `json:"review_id" db:"review_id"`
	ImageURL string `json:"image_url" db:"image_url"`
	Alt      string `json:"alt" db:"alt"`
	Position int    `json:"position" db:"position"`
}

// ReviewVote represents helpful/not helpful votes on reviews
type ReviewVote struct {
	ID        int64     `json:"id" db:"id"`
	ReviewID  int64     `json:"review_id" db:"review_id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	IsHelpful bool      `json:"is_helpful" db:"is_helpful"` // true for helpful, false for not helpful
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// ReviewResponse represents merchant responses to reviews
type ReviewResponse struct {
	ID        int64     `json:"id" db:"id"`
	ReviewID  int64     `json:"review_id" db:"review_id"`
	Response  string    `json:"response" db:"response"`
	CreatedBy int64     `json:"created_by" db:"created_by"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ProductRating represents aggregated product ratings
type ProductRating struct {
	ProductID     int64     `json:"product_id" db:"product_id"`
	AverageRating float64   `json:"average_rating" db:"average_rating"`
	TotalReviews  int       `json:"total_reviews" db:"total_reviews"`
	Rating1       int       `json:"rating_1" db:"rating_1"` // count of 1-star reviews
	Rating2       int       `json:"rating_2" db:"rating_2"` // count of 2-star reviews
	Rating3       int       `json:"rating_3" db:"rating_3"` // count of 3-star reviews
	Rating4       int       `json:"rating_4" db:"rating_4"` // count of 4-star reviews
	Rating5       int       `json:"rating_5" db:"rating_5"` // count of 5-star reviews
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// ReviewFilter represents filters for review queries
type ReviewFilter struct {
	ProductID  *int64     `json:"product_id"`
	UserID     *int64     `json:"user_id"`
	Rating     *int       `json:"rating"`
	IsVerified *bool      `json:"is_verified"`
	IsApproved *bool      `json:"is_approved"`
	DateFrom   *time.Time `json:"date_from"`
	DateTo     *time.Time `json:"date_to"`
	SortBy     string     `json:"sort_by"`    // newest, oldest, helpful, rating
	SortOrder  string     `json:"sort_order"` // asc, desc
}
