package dto

type ProductResponse struct {
	ID             uint   `json:"id"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	BasePrice      float64 `json:"base_price"`
	ImgURL         string `json:"img_url"`
	Stock          int    `json:"stock"`
	ProductType    string `json:"product_type"`
	CategoryID     uint   `json:"category_id,omitempty"`
	CategoryTitle  string `json:"category_title,omitempty"`
	TypeDrink      string `json:"type_drink,omitempty"`
	Volume         string `json:"volume,omitempty"`
	SizePizza      string `json:"size_pizza,omitempty"`
	ThicknessPizza string `json:"thickness_pizza,omitempty"`
}

type ProductFilter struct {
	Type       string `json:"type"` // pizza, drink, extra (+ legacy: snack, sauce, snacks, sauces, extras → extra)
	Search     string `json:"search"`
	Sort       string `json:"sort"` // price_asc, price_desc
	CategoryID uint   `json:"category_id"`
}

type ProductCategoryResponse struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type CreateProductCategoryInput struct {
	Title string `json:"title"`
}

type UpdateProductCategoryInput struct {
	Title *string `json:"title,omitempty"`
}

type CreatePizzaInput struct {
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Price            float64 `json:"price"`
	ImgURL           string  `json:"img_url"`
	SizePizzaID      uint    `json:"size_pizza_id"`
	ThicknessPizzaID uint    `json:"thickness_pizza_id"`
	Stock            *int    `json:"stock,omitempty"`
}

type UpdatePizzaInput struct {
	Title            *string  `json:"title,omitempty"`
	Description      *string  `json:"description,omitempty"`
	Price            *float64 `json:"price,omitempty"`
	ImgURL           *string  `json:"img_url,omitempty"`
	SizePizzaID      *uint    `json:"size_pizza_id,omitempty"`
	ThicknessPizzaID *uint    `json:"thickness_pizza_id,omitempty"`
	Stock            *int     `json:"stock,omitempty"`
}

type PizzaResponse struct {
	ID               uint    `json:"id"`
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	Price            float64 `json:"price"`
	ImgURL           string  `json:"img_url"`
	Stock            int     `json:"stock"`
	SizePizzaID      uint    `json:"size_pizza_id"`
	SizePizza        string  `json:"size_pizza"`
	ThicknessPizzaID uint   `json:"thickness_pizza_id"`
	ThicknessPizza   string  `json:"thickness_pizza"`
}

type CreateDrinkInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImgURL      string  `json:"img_url"`
	TypeDrinkID uint    `json:"type_drink_id"`
	VolumeID    uint    `json:"volume_id"`
	Stock       *int    `json:"stock,omitempty"`
}

type UpdateDrinkInput struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	ImgURL      *string  `json:"img_url,omitempty"`
	TypeDrinkID *uint    `json:"type_drink_id,omitempty"`
	VolumeID    *uint    `json:"volume_id,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
}

type DrinkResponse struct {
	ID          uint    `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImgURL      string  `json:"img_url"`
	Stock       int     `json:"stock"`
	TypeDrinkID uint    `json:"type_drink_id"`
	TypeDrink   string  `json:"type_drink"`
	VolumeID    uint    `json:"volume_id"`
	Volume      string  `json:"volume"`
}

type CreateExtraInput struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImgURL      string  `json:"img_url"`
	CategoryID  uint    `json:"category_id"`
	Stock       *int    `json:"stock,omitempty"`
}

type UpdateExtraInput struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	ImgURL      *string  `json:"img_url,omitempty"`
	CategoryID  *uint    `json:"category_id,omitempty"`
	Stock       *int     `json:"stock,omitempty"`
}

type ExtraResponse struct {
	ID            uint    `json:"id"`
	Title         string  `json:"title"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	ImgURL        string  `json:"img_url"`
	Stock         int     `json:"stock"`
	CategoryID    uint    `json:"category_id"`
	CategoryTitle string  `json:"category_title"`
}