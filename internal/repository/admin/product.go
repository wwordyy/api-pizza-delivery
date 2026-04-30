package adminRepo

import (
	"context"

	"gorm.io/gorm"

	"api-pizza-delivery/internal/models"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}


func (r *ProductRepository) GetAll(ctx context.Context, productType, search, sort string, categoryID uint) ([]models.Product, error) {
	var products []models.Product

	query := r.db.WithContext(ctx).Model(&models.Product{}).Preload("Category")

	if productType != "" {
		query = query.Where("product_type = ?", productType)
	}

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}

	// Сортировка по цене
	switch sort {
	case "price_asc":
		query = query.Order("base_price ASC")
	case "price_desc":
		query = query.Order("base_price DESC")
	default:
		query = query.Order("id ASC")
	}

	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.WithContext(ctx).Preload("Category").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) GetAllByType(ctx context.Context, productType string) ([]models.Product, error) {
	var products []models.Product
	err := r.db.WithContext(ctx).
		Preload("Category").
		Where("product_type = ?", productType).
		Order("id ASC").
		Find(&products).Error
	return products, err
}


func (r *ProductRepository) Create(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductRepository) Update(ctx context.Context, p *models.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *ProductRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Product{}, id).Error
}

func (r *ProductRepository) CreatePizza(ctx context.Context, p *models.Pizza) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProductRepository) GetPizzaByProductID(ctx context.Context, productID uint) (*models.Pizza, error) {
	var pizza models.Pizza
	err := r.db.WithContext(ctx).
		Preload("SizePizza").
		Preload("ThicknessPizza").
		Where("product_id = ?", productID).
		First(&pizza).Error
	if err != nil {
		return nil, err
	}
	return &pizza, nil
}

func (r *ProductRepository) UpdatePizza(ctx context.Context, p *models.Pizza) error {
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *ProductRepository) CreateDrink(ctx context.Context, d *models.Drink) error {
	return r.db.WithContext(ctx).Create(d).Error
}

func (r *ProductRepository) GetDrinkByProductID(ctx context.Context, productID uint) (*models.Drink, error) {
	var drink models.Drink
	err := r.db.WithContext(ctx).
		Preload("TypeDrink").
		Preload("VolumeDrink").
		Where("product_id = ?", productID).
		First(&drink).Error
	if err != nil {
		return nil, err
	}
	return &drink, nil
}

func (r *ProductRepository) UpdateDrink(ctx context.Context, d *models.Drink) error {
	return r.db.WithContext(ctx).Save(d).Error
}

func (r *ProductRepository) GetSizePizzas(ctx context.Context) ([]models.SizePizza, error) {
	var sizes []models.SizePizza
	err := r.db.WithContext(ctx).Find(&sizes).Error
	return sizes, err
}

func (r *ProductRepository) GetThicknessPizzas(ctx context.Context) ([]models.ThicknessPizza, error) {
	var thicknesses []models.ThicknessPizza
	err := r.db.WithContext(ctx).Find(&thicknesses).Error
	return thicknesses, err
}

func (r *ProductRepository) GetTypeDrinks(ctx context.Context) ([]models.TypeDrink, error) {
	var types []models.TypeDrink
	err := r.db.WithContext(ctx).Find(&types).Error
	return types, err
}

func (r *ProductRepository) GetVolumeDrinks(ctx context.Context) ([]models.VolumeDrink, error) {
	var volumes []models.VolumeDrink
	err := r.db.WithContext(ctx).Find(&volumes).Error
	return volumes, err
}
