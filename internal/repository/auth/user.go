

package authRepo

import (
	"api-pizza-delivery/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, u *models.User) error {
	return r.db.WithContext(ctx).Create(u).Error	
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("email = ?", email).
		First(&u).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("пользователь не найден")
	}
	
	
	return &u, err
}

func (r *UserRepository) FindById(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	
	err := r.db.WithContext(ctx).
		Preload("Role").
		Where("id = ?", id).
		First(&u).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("пользователь не найден")
	}
	
	
	return &u, err
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	
	var count int64
	
	err := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email).
		Count(&count).Error
	
	return count > 0, err
}

func (r *UserRepository) UpdatePassword(ctx context.Context, 
										email, hashedPswd string) error {
	
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("email = ?", email).
		Update("password", hashedPswd).Error
}