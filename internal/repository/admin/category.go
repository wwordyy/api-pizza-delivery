package adminRepo

import (
	"context"

	"api-pizza-delivery/internal/models"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) List(ctx context.Context) ([]models.ProductCategory, error) {
	var list []models.ProductCategory
	err := r.db.WithContext(ctx).Order("id ASC").Find(&list).Error
	return list, err
}

func (r *CategoryRepository) GetByID(ctx context.Context, id uint) (*models.ProductCategory, error) {
	var c models.ProductCategory
	err := r.db.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Create(ctx context.Context, c *models.ProductCategory) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *CategoryRepository) Update(ctx context.Context, c *models.ProductCategory) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *CategoryRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.ProductCategory{}, id).Error
}

func (r *CategoryRepository) CountProducts(ctx context.Context, categoryID uint) (int64, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&models.Product{}).Where("category_id = ?", categoryID).Count(&n).Error
	return n, err
}
