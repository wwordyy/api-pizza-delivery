package adminRepo

import (
	"api-pizza-delivery/internal/models"
	"context"
	"time"

	"gorm.io/gorm"
)

type AnalyticsRepository struct {
	db *gorm.DB
}


func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	return &AnalyticsRepository{db: db}
}

type AnalyticsFilter struct {
	From            *time.Time
	To              *time.Time
	UserID          *uint
	PaymentMethodID *uint
	Sort            string
}

func (r *AnalyticsRepository) GetOrders(ctx context.Context, f AnalyticsFilter) ([]models.Order, error) {
	var orders []models.Order

	query := r.db.WithContext(ctx).
		Preload("PaymentMethod").
		Model(&models.Order{})

	if f.From != nil {
		query = query.Where("date_of_order >= ?", f.From)
	}
	if f.To != nil {
		query = query.Where("date_of_order <= ?", f.To)
	}
	if f.UserID != nil {
		query = query.Where("user_id = ?", f.UserID)
	}
	if f.PaymentMethodID != nil {
		query = query.Where("payment_method_id = ?", f.PaymentMethodID)
	}

	switch f.Sort {
	case "date_asc":
		query = query.Order("date_of_order ASC")
	case "date_desc":
		query = query.Order("date_of_order DESC")
	default:
		query = query.Order("created_at DESC")
	}

	return orders, query.Find(&orders).Error
}

type AnalyticsSummary struct {
	TotalSum      float64
	TotalCount    int64
	AverageOrder  float64
}

func (r *AnalyticsRepository) GetSummary(ctx context.Context, f AnalyticsFilter) (*AnalyticsSummary, error) {
	query := r.db.WithContext(ctx).Model(&models.Order{})

	if f.From != nil {
		query = query.Where("date_of_order >= ?", f.From)
	}
	if f.To != nil {
		query = query.Where("date_of_order <= ?", f.To)
	}
	if f.UserID != nil {
		query = query.Where("user_id = ?", f.UserID)
	}
	if f.PaymentMethodID != nil {
		query = query.Where("payment_method_id = ?", f.PaymentMethodID)
	}

	var summary AnalyticsSummary
	err := query.Select(
		"COALESCE(SUM(sum_order), 0) as total_sum, COUNT(*) as total_count, COALESCE(AVG(sum_order), 0) as average_order",
	).Scan(&summary).Error

	return &summary, err
}
