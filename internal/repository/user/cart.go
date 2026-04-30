package userRepo

import (
	"api-pizza-delivery/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) GetOrCreate(ctx context.Context, userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&cart).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		cart = models.Cart{UserID: userID}
		if err := r.db.WithContext(ctx).Create(&cart).Error; err != nil {
			return nil, err
		}
		return &cart, nil
	}
	return &cart, err
}

func (r *CartRepository) GetWithItems(ctx context.Context, userID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.WithContext(ctx).
		Preload("CartItems.Product").
		Where("user_id = ?", userID).
		First(&cart).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("корзина не найдена")
	}
	return &cart, err
}

func (r *CartRepository) AddItem(ctx context.Context, cartID, productID uint, quantity int) error {
	var item models.CartItem

	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		First(&item).Error

	if err == nil {
		return r.db.WithContext(ctx).
			Model(&item).
			Update("quantity", item.Quantity+quantity).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		newItem := models.CartItem{
			CartID:    cartID,
			ProductID: productID,
			Quantity:  quantity,
		}
		return r.db.WithContext(ctx).Create(&newItem).Error
	}

	return err
}

func (r *CartRepository) UpdateItem(ctx context.Context, itemID uint, quantity int) error {
	return r.db.WithContext(ctx).
		Model(&models.CartItem{}).
		Where("id = ?", itemID).
		Update("quantity", quantity).Error
}


func (r *CartRepository) DeleteItem(ctx context.Context, itemID uint) error {
	return r.db.WithContext(ctx).
		Delete(&models.CartItem{}, itemID).Error
}

func (r *CartRepository) GetItem(ctx context.Context, itemID uint) (*models.CartItem, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).
		Preload("Product").
		First(&item, itemID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("позиция в корзине не найдена")
	}
	return &item, err
}

func (r *CartRepository) LineQuantity(ctx context.Context, cartID, productID uint) (int, error) {
	var item models.CartItem
	err := r.db.WithContext(ctx).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		First(&item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return item.Quantity, nil
}

func (r *CartRepository) GetProductStock(ctx context.Context, productID uint) (int, error) {
	var p models.Product
	err := r.db.WithContext(ctx).Select("stock").First(&p, productID).Error
	if err != nil {
		return 0, err
	}
	return p.Stock, nil
}
