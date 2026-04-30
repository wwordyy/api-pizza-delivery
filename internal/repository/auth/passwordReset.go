package authRepo

import (
	"api-pizza-delivery/internal/models"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type PasswordResetRepository struct {
	db *gorm.DB
}

func NewPasswordResetRepository(db *gorm.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}


func (r *PasswordResetRepository) Create(ctx context.Context, email, code string,
											expireAt time.Time) error {
	
	reset := models.PasswordResetCode{
		Email:    email,
		Code:     code,
		ExpiresAt: expireAt,
	}
	
	return r.db.WithContext(ctx).Create(&reset).Error
	
}

func (r *PasswordResetRepository) FindValid(ctx context.Context, code, email string) (*models.PasswordResetCode, error) {
	
	var reset models.PasswordResetCode

	err := r.db.WithContext(ctx).
		Where("email = ? AND code = ? AND is_used = ? AND expires_at > ?", email, code, false, time.Now()).
		First(&reset).Error
	
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("код сброса пароля не найден")
	}
	
	return &reset, err
}

func (r *PasswordResetRepository) MarkUsed(ctx context.Context, id uint) error {
	
	return r.db.WithContext(ctx).
		Model(&models.PasswordResetCode{}).
		Where("id = ?", id).
		Update("is_used", true).Error
}
