package adminRepo

import (
	"api-pizza-delivery/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type ReviewRepository struct {
	db *gorm.DB
}


func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

func (r *ReviewRepository) GetAll(ctx context.Context) ([]models.Review, error) {
	
	var reviews []models.Review
	
	err := r.db.WithContext(ctx).
		Preload("User").
		Preload("User.Role").
		Preload("Product").
		Find(&reviews).Error
	return reviews, err
}

func (r *ReviewRepository) GetByProductID(ctx context.Context, productID uint) ([]models.Review, error) {
	var reviews []models.Review
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Find(&reviews).Error
	return reviews, err
}


func (r *ReviewRepository) Create(ctx context.Context, review *models.Review) error {
	return r.db.WithContext(ctx).Create(review).Error
}

func (r *ReviewRepository) GetByID(ctx context.Context, id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.WithContext(ctx).First(&review, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	return &review, err
}

func (r *ReviewRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Review{}, id).Error
}