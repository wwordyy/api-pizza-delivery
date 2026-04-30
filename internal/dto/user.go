package dto

import "time"

type ProfileResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type UpdateProfileInput struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

type ChangePasswordInput struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

type CartResponse struct {
	Items []CartItemResponse `json:"items"`
	Total float64            `json:"total"`
}

type CartItemResponse struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Title     string  `json:"title"`
	Price     float64 `json:"price"`
	Quantity  int     `json:"quantity"`
	Stock     int     `json:"stock"`
	Subtotal  float64 `json:"subtotal"`
}

type AddToCartInput struct {
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"`
}

type UpdateCartItemInput struct {
	Quantity int `json:"quantity"`
}

// Заказы
type CreateOrderInput struct {
	DeliveryDate    string `json:"delivery_date"`
	Addresses       string `json:"addresses"`
	PaymentMethodID uint   `json:"payment_method_id"`
	BankCard        string `json:"bank_card"`
}

type OrderResponse struct {
	ID            uint       `json:"id"`
	UserID        uint       `json:"user_id"`
	DateOfOrder   time.Time  `json:"date_of_order"`
	DeliveryDate  *time.Time `json:"delivery_date,omitempty"`
	SumOrder      float64    `json:"sum_order"`
	Addresses     string     `json:"addresses"`
	PaymentMethod string     `json:"payment_method"`
	StatusOrder   string     `json:"status_order"`
	CreatedAt     time.Time  `json:"created_at"`
}

type DeliveryAddressResponse struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Label     string    `json:"label"`
	Line      string    `json:"line"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateDeliveryAddressInput struct {
	Label     string `json:"label"`
	Line      string `json:"line"`
	IsDefault bool   `json:"is_default"`
}

type UpdateDeliveryAddressInput struct {
	Label     *string `json:"label"`
	Line      *string `json:"line"`
	IsDefault *bool   `json:"is_default"`
}

// Избранное
type FavoriteResponse struct {
	ID          uint    `json:"id"`
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Price       float64 `json:"price"`
	ImgURL      string  `json:"img_url"`
}

type AddFavoriteInput struct {
	ProductID uint `json:"product_id"`
}

type ReviewResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	ProductID   uint      `json:"product_id"`
	Rating      int       `json:"rating"`
	Comment     string    `json:"comment"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateReviewInput struct {
	Rating  int    `json:"rating"`
	Comment string `json:"comment"`
}
