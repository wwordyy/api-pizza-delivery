package models

import "time"

type ProductCategory struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"uniqueIndex;not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"not null"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Phone     string    `json:"phone"`
	Password  string    `json:"-" gorm:"not null"`
	RoleID    uint      `json:"role_id" gorm:"default:1"`
	Role      Role      `json:"role" gorm:"foreignKey:RoleID"`
	CreatedAt time.Time `json:"created_at"`
}

type Role struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"uniqueIndex;not null"`
}

type Product struct {
	ID          uint             `json:"id" gorm:"primaryKey"`
	Title       string           `json:"title" gorm:"not null"`
	Description string           `json:"description"`
	BasePrice   float64          `json:"base_price" gorm:"not null"`
	ImgURL      string           `json:"img_url"`
	ProductType string           `json:"product_type" gorm:"not null"`
	Stock       int              `json:"stock" gorm:"not null;default:5"`
	CategoryID  *uint            `json:"category_id,omitempty" gorm:"index"`
	Category    *ProductCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
}

type Pizza struct {
	ID               uint           `json:"id" gorm:"primaryKey"`
	ProductID        uint           `json:"product_id"`
	Product          Product        `json:"product" gorm:"foreignKey:ProductID"`
	SizePizzaID      uint           `json:"size_pizza_id"`
	SizePizza        SizePizza      `json:"size_pizza" gorm:"foreignKey:SizePizzaID"`
	ThicknessPizzaID uint           `json:"thickness_pizza_id"`
	ThicknessPizza   ThicknessPizza `json:"thickness_pizza" gorm:"foreignKey:ThicknessPizzaID"`
}

type SizePizza struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"uniqueIndex;not null"`
}

type ThicknessPizza struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"uniqueIndex;not null"`
}

type Drink struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	ProductID   uint       `json:"product_id"`
	Product     Product    `json:"product" gorm:"foreignKey:ProductID"`
	TypeDrinkID uint       `json:"type_drink_id"`
	TypeDrink   TypeDrink  `json:"type_drink" gorm:"foreignKey:TypeDrinkID"`
	VolumeID    uint       `json:"volume_id"`
	VolumeDrink VolumeDrink `json:"volume_drink" gorm:"foreignKey:VolumeID"`
}

type TypeDrink struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"uniqueIndex;not null"`
}

type VolumeDrink struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"uniqueIndex;not null"`
}

type Cart struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    uint       `json:"user_id" gorm:"uniqueIndex"`
	CartItems []CartItem `json:"cart_items" gorm:"foreignKey:CartID"`
}

type CartItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	CartID    uint    `json:"cart_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity" gorm:"default:1"`
}

type Order struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	UserID          uint          `json:"user_id"`
	DateOfOrder     time.Time     `json:"date_of_order"`
	DeliveryDate    *time.Time    `json:"delivery_date"`
	SumOrder        float64       `json:"sum_order" gorm:"not null"`
	Addresses       string        `json:"addresses"`
	PaymentMethodID uint          `json:"payment_method_id"`
	PaymentMethod   PaymentMethod `json:"payment_method" gorm:"foreignKey:PaymentMethodID"`
	StatusOrderID   uint          `json:"status_order_id" gorm:"default:1"`
	StatusOrder     StatusOrder   `json:"status_order" gorm:"foreignKey:StatusOrderID"`
	CreatedAt       time.Time     `json:"created_at"`
	OrderItems      []OrderItem   `json:"order_items,omitempty" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	OrderID   uint    `json:"order_id" gorm:"index;not null"`
	ProductID uint    `json:"product_id" gorm:"not null"`
	Product   Product `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	UnitPrice float64 `json:"unit_price" gorm:"not null"`
}

type PaymentMethod struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"uniqueIndex;not null"`
}

type StatusOrder struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"uniqueIndex;not null"`
}

type DeliveryAddress struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uint      `json:"user_id" gorm:"index;not null"`
	Label     string    `json:"label"`
	Line      string    `json:"line" gorm:"not null"`
	IsDefault bool      `json:"is_default" gorm:"default:false"`
}

type Favorite struct {
	ID        uint    `json:"id" gorm:"primaryKey"`
	UserID    uint    `json:"user_id"`
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
}

type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserID"`
	ProductID uint      `json:"product_id"`
	Product   Product   `json:"product" gorm:"foreignKey:ProductID"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type PasswordResetCode struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"not null"`
	Code      string    `json:"code" gorm:"not null"`
	IsUsed    bool      `json:"is_used" gorm:"default:false"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}