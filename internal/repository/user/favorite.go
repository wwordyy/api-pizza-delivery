package userRepo

import (
	"api-pizza-delivery/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{db: db}
}

func (r *FavoriteRepository) GetAll(ctx context.Context, userID uint) ([]models.Favorite, error) {
	var favorites []models.Favorite
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("user_id = ?", userID).
		Find(&favorites).Error
	return favorites, err
}

func (r *FavoriteRepository) Add(ctx context.Context, userID, productID uint) error {
	favorite := models.Favorite{
		UserID:    userID,
		ProductID: productID,
	}
	return r.db.WithContext(ctx).Create(&favorite).Error
}

func (r *FavoriteRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Favorite{}, id).Error
}

func (r *FavoriteRepository) GetByID(ctx context.Context, id uint) (*models.Favorite, error) {
	var favorite models.Favorite
	err := r.db.WithContext(ctx).First(&favorite, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("избранное не найдено")
	}
	return &favorite, err
}

func (r *FavoriteRepository) Exists(ctx context.Context, userID, productID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Favorite{}).
		Where("user_id = ? AND product_id = ?", userID, productID).
		Count(&count).Error
	return count > 0, err
}