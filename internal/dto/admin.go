package dto

import "time"

type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}




type AnalyticsFilter struct {
	From            string `json:"from"` // формат 2024-01-01
	To              string `json:"to"`   // формат 2024-12-31
	UserID          uint   `json:"user_id"`
	PaymentMethodID uint   `json:"payment_method_id"`
	Sort            string `json:"sort"` // date_asc, date_desc
}



type AnalyticsResponse struct {
	Summary AnalyticsSummaryDTO `json:"summary"`
	Orders  []OrderResponse     `json:"orders"`
}

type AnalyticsSummaryDTO struct {
	TotalSum     float64 `json:"total_sum"`
	TotalCount   int64   `json:"total_count"`
	AverageOrder float64 `json:"average_order"`
}

type UploadImageResponse struct {
	URL string `json:"url"`
}

type AdminReviewResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	ProductID   uint      `json:"product_id"`
	ProductName string    `json:"product_name"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}
