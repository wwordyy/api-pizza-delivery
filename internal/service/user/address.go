package userService

import (
	"api-pizza-delivery/internal/dto"
	"api-pizza-delivery/internal/models"
	user "api-pizza-delivery/internal/repository/user"
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrAddressNotFound  = errors.New("адрес доставки не найден")
	ErrAddressNotYours  = errors.New("этот адрес не принадлежит пользователю")
	ErrAddressLineEmpty = errors.New("адрес не может быть пустым")
)

type AddressService struct {
	repo *user.DeliveryAddressRepository
}

func NewAddressService(repo *user.DeliveryAddressRepository) *AddressService {
	return &AddressService{repo: repo}
}

func (s *AddressService) List(ctx context.Context, userID uint) ([]dto.DeliveryAddressResponse, error) {
	rows, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]dto.DeliveryAddressResponse, len(rows))
	for i, a := range rows {
		out[i] = toAddressResponse(a)
	}
	return out, nil
}

func (s *AddressService) Create(ctx context.Context, userID uint, in dto.CreateDeliveryAddressInput) (*dto.DeliveryAddressResponse, error) {
	if in.Line == "" {
		return nil, errors.New("укажите адрес (поле line)")
	}
	addr := &models.DeliveryAddress{
		UserID:    userID,
		Label:     in.Label,
		Line:      in.Line,
		IsDefault: in.IsDefault,
	}
	if err := s.repo.Create(ctx, addr); err != nil {
		return nil, err
	}
	created, err := s.repo.GetByID(ctx, addr.ID)
	if err != nil {
		return nil, err
	}
	resp := toAddressResponse(*created)
	return &resp, nil
}

func (s *AddressService) Update(ctx context.Context, userID uint, id uint, in dto.UpdateDeliveryAddressInput) (*dto.DeliveryAddressResponse, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, user.ErrDeliveryAddressNotFound) {
			return nil, ErrAddressNotFound
		}
		return nil, err
	}
	if existing.UserID != userID {
		return nil, ErrAddressNotYours
	}
	if in.Line != nil && *in.Line == "" {
		return nil, ErrAddressLineEmpty
	}
	updated := *existing
	if in.Label != nil {
		updated.Label = *in.Label
	}
	if in.Line != nil {
		updated.Line = *in.Line
	}
	if in.IsDefault != nil {
		updated.IsDefault = *in.IsDefault
	}
	if err := s.repo.Update(ctx, userID, &updated); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrAddressNotFound
		}
		return nil, err
	}
	refreshed, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := toAddressResponse(*refreshed)
	return &resp, nil
}

func (s *AddressService) Delete(ctx context.Context, userID uint, id uint) error {
	err := s.repo.Delete(ctx, userID, id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrAddressNotFound
	}
	return err
}

func toAddressResponse(a models.DeliveryAddress) dto.DeliveryAddressResponse {
	return dto.DeliveryAddressResponse{
		ID:        a.ID,
		UserID:    a.UserID,
		Label:     a.Label,
		Line:      a.Line,
		IsDefault: a.IsDefault,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}
