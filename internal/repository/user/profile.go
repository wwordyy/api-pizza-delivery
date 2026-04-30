package userRepo

import (
	"api-pizza-delivery/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		First(&user, id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("пользователь не найден")
	}
	return &user, err
}

func (r *ProfileRepository) Update(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).
		Model(user).
		Updates(map[string]interface{}{
			"username": user.Username,
			"phone":    user.Phone,
		}).Error
}

func (r *ProfileRepository) UpdatePassword(ctx context.Context, id uint, hashedPassword string) error {
	return r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Update("password", hashedPassword).Error
}

func (r *ProfileRepository) GetPasswordByID(ctx context.Context, id uint) (string, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Select("password").
		First(&user, id).Error
	if err != nil {
		return "", err
	}
	return user.Password, nil
}