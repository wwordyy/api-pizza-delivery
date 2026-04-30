package userService

import (
	"api-pizza-delivery/internal/dto"
	user "api-pizza-delivery/internal/repository/user"
	"context"
	"errors"

	"gorm.io/gorm"
)

var (
	ErrCartItemNotFound    = errors.New("позиция в корзине не найдена")
	ErrInvalidQuantity     = errors.New("количество должно быть больше 0")
	ErrNotYourItem         = errors.New("эта позиция не принадлежит вашей корзине")
	ErrInsufficientStock   = errors.New("недостаточно товара на складе")
	ErrProductNotFoundCart = errors.New("товар не найден")
)

type CartService struct {
	cartRepo *user.CartRepository
}

func NewCartService(cartRepo *user.CartRepository) *CartService {
	return &CartService{cartRepo: cartRepo}
}

func (s *CartService) GetCart(ctx context.Context, userID uint) (*dto.CartResponse, error) {
	cart, err := s.cartRepo.GetWithItems(ctx, userID)
	if err != nil {
		// Если корзины нет — возвращаем пустую
		return &dto.CartResponse{Items: []dto.CartItemResponse{}}, nil
	}

	response := &dto.CartResponse{}
	var total float64

	for _, item := range cart.CartItems {
		response.Items = append(response.Items, dto.CartItemResponse{
			ID:        item.ID,
			ProductID: item.ProductID,
			Title:     item.Product.Title,
			Price:     item.Product.BasePrice,
			Quantity:  item.Quantity,
			Stock:     item.Product.Stock,
			Subtotal:  item.Product.BasePrice * float64(item.Quantity),
		})
		total += item.Product.BasePrice * float64(item.Quantity)
	}

	response.Total = total
	return response, nil
}

func (s *CartService) AddItem(ctx context.Context, userID, productID uint, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	cart, err := s.cartRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return err
	}

	stock, err := s.cartRepo.GetProductStock(ctx, productID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFoundCart
		}
		return err
	}
	current, err := s.cartRepo.LineQuantity(ctx, cart.ID, productID)
	if err != nil {
		return err
	}
	if current+quantity > stock {
		return ErrInsufficientStock
	}

	return s.cartRepo.AddItem(ctx, cart.ID, productID, quantity)
}

func (s *CartService) UpdateItem(ctx context.Context, userID, itemID uint, quantity int) error {
	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	item, err := s.cartRepo.GetItem(ctx, itemID)
	if err != nil {
		return ErrCartItemNotFound
	}

	// Проверяем что позиция принадлежит корзине пользователя
	cart, err := s.cartRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return err
	}
	if item.CartID != cart.ID {
		return ErrNotYourItem
	}

	stock, err := s.cartRepo.GetProductStock(ctx, item.ProductID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFoundCart
		}
		return err
	}
	if quantity > stock {
		return ErrInsufficientStock
	}

	return s.cartRepo.UpdateItem(ctx, itemID, quantity)
}

func (s *CartService) DeleteItem(ctx context.Context, userID, itemID uint) error {
	item, err := s.cartRepo.GetItem(ctx, itemID)
	if err != nil {
		return ErrCartItemNotFound
	}

	cart, err := s.cartRepo.GetOrCreate(ctx, userID)
	if err != nil {
		return err
	}
	if item.CartID != cart.ID {
		return ErrNotYourItem
	}

	return s.cartRepo.DeleteItem(ctx, itemID)
}
