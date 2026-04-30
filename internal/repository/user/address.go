package userRepo

import (
	"api-pizza-delivery/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

var ErrDeliveryAddressNotFound = errors.New("адрес доставки не найден")

type DeliveryAddressRepository struct {
	db *gorm.DB
}

func NewDeliveryAddressRepository(db *gorm.DB) *DeliveryAddressRepository {
	return &DeliveryAddressRepository{db: db}
}

func (r *DeliveryAddressRepository) ListByUserID(ctx context.Context, userID uint) ([]models.DeliveryAddress, error) {
	var rows []models.DeliveryAddress
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("is_default DESC, created_at DESC").
		Find(&rows).Error
	return rows, err
}

func (r *DeliveryAddressRepository) GetByID(ctx context.Context, id uint) (*models.DeliveryAddress, error) {
	var row models.DeliveryAddress
	err := r.db.WithContext(ctx).First(&row, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrDeliveryAddressNotFound
	}
	return &row, err
}

func (r *DeliveryAddressRepository) Create(ctx context.Context, addr *models.DeliveryAddress) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if addr.IsDefault {
			if err := tx.Model(&models.DeliveryAddress{}).
				Where("user_id = ?", addr.UserID).
				Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(addr).Error
	})
}

func (r *DeliveryAddressRepository) Update(ctx context.Context, userID uint, addr *models.DeliveryAddress) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var existing models.DeliveryAddress
		if err := tx.Where("id = ? AND user_id = ?", addr.ID, userID).First(&existing).Error; err != nil {
			return err
		}
		if addr.IsDefault {
			if err := tx.Model(&models.DeliveryAddress{}).
				Where("user_id = ? AND id != ?", userID, addr.ID).
				Update("is_default", false).Error; err != nil {
				return err
			}
		}
		return tx.Model(&existing).Updates(map[string]interface{}{
			"label":      addr.Label,
			"line":       addr.Line,
			"is_default": addr.IsDefault,
		}).Error
	})
}

func (r *DeliveryAddressRepository) Delete(ctx context.Context, userID uint, id uint) error {
	res := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&models.DeliveryAddress{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
