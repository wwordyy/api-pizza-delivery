package adminRepo

import (
	"api-pizza-delivery/internal/models"
	"context"

	"gorm.io/gorm"
)

type AdminUserRepository struct {
	db *gorm.DB
}

func NewAdminUserRepository(db *gorm.DB) *AdminUserRepository {
	return &AdminUserRepository{db: db}
}

func (s *AdminUserRepository) GetAll(ctx context.Context) ([]models.User, error) {

	var users []models.User
	err := s.db.WithContext(ctx).
		Preload("Role").
		Find(&users).Error
	return users, err
	
}
