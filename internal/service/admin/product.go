package adminService

import (
	adminRepo "api-pizza-delivery/internal/repository/admin"
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"

	"api-pizza-delivery/internal/dto"
	"api-pizza-delivery/internal/models"
)

var (
	ErrProductNotFound  = errors.New("товар не найден")
	ErrNotExtraProduct  = errors.New("это не дополнительный товар")
	ErrInvalidPrice     = errors.New("неверная цена")
	ErrInvalidTitle     = errors.New("неверное название")
	ErrInvalidDesc      = errors.New("неверное описание")
	ErrInvalidStock     = errors.New("неверное количество на складе")
	ErrNoFieldsToUpdate = errors.New("нет полей для обновления")
)

type ProductService struct {
	productRepo  *adminRepo.ProductRepository
	categoryRepo *adminRepo.CategoryRepository
}

func NewProductService(productRepo *adminRepo.ProductRepository, categoryRepo *adminRepo.CategoryRepository) *ProductService {
	return &ProductService{productRepo: productRepo, categoryRepo: categoryRepo}
}

func normalizeProductListType(t string) string {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "extras", "extra", "snack", "snacks", "sauce", "sauces":
		return "extra"
	default:
		return strings.TrimSpace(t)
	}
}


func (s *ProductService) GetAll(ctx context.Context, filter dto.ProductFilter) ([]dto.ProductResponse, error) {
	pt := normalizeProductListType(filter.Type)
	products, err := s.productRepo.GetAll(ctx, pt, filter.Search, filter.Sort, filter.CategoryID)
	if err != nil {
		return nil, err
	}

	out := make([]dto.ProductResponse, 0, len(products))
	for i := range products {
		p := products[i]
		resp, err := s.toProductResponse(ctx, &p)
		if err != nil {
			return nil, err
		}
		out = append(out, *resp)
	}
	return out, nil
}

func (s *ProductService) GetAllPizzas(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := s.productRepo.GetAllByType(ctx, "pizza")
	if err != nil {
		return nil, err
	}

	out := make([]dto.ProductResponse, 0, len(products))
	for i := range products {
		resp, err := s.toProductResponse(ctx, &products[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *resp)
	}
	return out, nil
}

func (s *ProductService) GetAllDrinks(ctx context.Context) ([]dto.ProductResponse, error) {
	products, err := s.productRepo.GetAllByType(ctx, "drink")
	if err != nil {
		return nil, err
	}

	out := make([]dto.ProductResponse, 0, len(products))
	for i := range products {
		resp, err := s.toProductResponse(ctx, &products[i])
		if err != nil {
			return nil, err
		}
		out = append(out, *resp)
	}
	return out, nil
}

func (s *ProductService) GetAllExtras(ctx context.Context) ([]dto.ExtraResponse, error) {
	products, err := s.productRepo.GetAllByType(ctx, "extra")
	if err != nil {
		return nil, err
	}
	out := make([]dto.ExtraResponse, 0, len(products))
	for i := range products {
		out = append(out, *toExtraResponse(&products[i]))
	}
	return out, nil
}

func (s *ProductService) GetByID(ctx context.Context, id uint) (*dto.ProductResponse, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	return s.toProductResponse(ctx, product)
}


func (s *ProductService) CreatePizza(ctx context.Context, input dto.CreatePizzaInput) (*dto.PizzaResponse, error) {
	if err := validatePizzaInput(input.Title, input.Description, input.Price); err != nil {
		return nil, err
	}
	stock, err := stockForCreate(input.Stock)
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		Title:       input.Title,
		Description: input.Description,
		BasePrice:   input.Price,
		ImgURL:      input.ImgURL,
		ProductType: "pizza",
		Stock:       stock,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	pizza := &models.Pizza{
		ProductID:        product.ID,
		SizePizzaID:      input.SizePizzaID,
		ThicknessPizzaID: input.ThicknessPizzaID,
	}

	if err := s.productRepo.CreatePizza(ctx, pizza); err != nil {
		return nil, err
	}

	pizza, err = s.productRepo.GetPizzaByProductID(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	return toPizzaResponse(product, pizza), nil
}


func (s *ProductService) UpdatePizza(ctx context.Context, id uint, input dto.UpdatePizzaInput) (*dto.PizzaResponse, error) {
	if input.Title == nil && input.Description == nil && input.Price == nil && input.ImgURL == nil &&
		input.SizePizzaID == nil && input.ThicknessPizzaID == nil && input.Stock == nil {
		return nil, ErrNoFieldsToUpdate
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	needsProductValidation := input.Title != nil || input.Description != nil || input.Price != nil
	if needsProductValidation {
		title := product.Title
		description := product.Description
		price := product.BasePrice
		if input.Title != nil {
			title = *input.Title
		}
		if input.Description != nil {
			description = *input.Description
		}
		if input.Price != nil {
			price = *input.Price
		}

		if err := validatePizzaInput(title, description, price); err != nil {
			return nil, err
		}
	}

	if input.Title != nil {
		product.Title = *input.Title
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil {
		product.BasePrice = *input.Price
	}
	if input.ImgURL != nil {
		product.ImgURL = *input.ImgURL
	}
	if input.Stock != nil {
		if *input.Stock < 0 {
			return nil, ErrInvalidStock
		}
		product.Stock = *input.Stock
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	pizza, err := s.productRepo.GetPizzaByProductID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.SizePizzaID != nil {
		pizza.SizePizzaID = *input.SizePizzaID
	}
	if input.ThicknessPizzaID != nil {
		pizza.ThicknessPizzaID = *input.ThicknessPizzaID
	}

	if err := s.productRepo.UpdatePizza(ctx, pizza); err != nil {
		return nil, err
	}

	pizza, err = s.productRepo.GetPizzaByProductID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toPizzaResponse(product, pizza), nil
}

func (s *ProductService) DeletePizza(ctx context.Context, id uint) error {
	_, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}
	return s.productRepo.Delete(ctx, id)
}


func (s *ProductService) CreateDrink(ctx context.Context, input dto.CreateDrinkInput) (*dto.DrinkResponse, error) {
	if err := validateDrinkInput(input.Title, input.Description, input.Price); err != nil {
		return nil, err
	}
	stock, err := stockForCreate(input.Stock)
	if err != nil {
		return nil, err
	}

	product := &models.Product{
		Title:       input.Title,
		Description: input.Description,
		BasePrice:   input.Price,
		ImgURL:      input.ImgURL,
		ProductType: "drink",
		Stock:       stock,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	drink := &models.Drink{
		ProductID:   product.ID,
		TypeDrinkID: input.TypeDrinkID,
		VolumeID:    input.VolumeID,
	}

	if err := s.productRepo.CreateDrink(ctx, drink); err != nil {
		return nil, err
	}

	drink, err = s.productRepo.GetDrinkByProductID(ctx, product.ID)
	if err != nil {
		return nil, err
	}

	return toDrinkResponse(product, drink), nil
}

func (s *ProductService) UpdateDrink(ctx context.Context, id uint, input dto.UpdateDrinkInput) (*dto.DrinkResponse, error) {
	if input.Title == nil && input.Description == nil && input.Price == nil && input.ImgURL == nil &&
		input.TypeDrinkID == nil && input.VolumeID == nil && input.Stock == nil {
		return nil, ErrNoFieldsToUpdate
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}

	needsProductValidation := input.Title != nil || input.Description != nil || input.Price != nil
	if needsProductValidation {
		title := product.Title
		description := product.Description
		price := product.BasePrice
		if input.Title != nil {
			title = *input.Title
		}
		if input.Description != nil {
			description = *input.Description
		}
		if input.Price != nil {
			price = *input.Price
		}

		if err := validateDrinkInput(title, description, price); err != nil {
			return nil, err
		}
	}

	if input.Title != nil {
		product.Title = *input.Title
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil {
		product.BasePrice = *input.Price
	}
	if input.ImgURL != nil {
		product.ImgURL = *input.ImgURL
	}
	if input.Stock != nil {
		if *input.Stock < 0 {
			return nil, ErrInvalidStock
		}
		product.Stock = *input.Stock
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	drink, err := s.productRepo.GetDrinkByProductID(ctx, id)
	if err != nil {
		return nil, err
	}

	if input.TypeDrinkID != nil {
		drink.TypeDrinkID = *input.TypeDrinkID
	}
	if input.VolumeID != nil {
		drink.VolumeID = *input.VolumeID
	}

	if err := s.productRepo.UpdateDrink(ctx, drink); err != nil {
		return nil, err
	}

	drink, err = s.productRepo.GetDrinkByProductID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toDrinkResponse(product, drink), nil
}

func (s *ProductService) DeleteDrink(ctx context.Context, id uint) error {
	_, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}
	return s.productRepo.Delete(ctx, id)
}

func (s *ProductService) CreateExtra(ctx context.Context, input dto.CreateExtraInput) (*dto.ExtraResponse, error) {
	if err := validateSimpleProduct(input.Title, input.Description, input.Price); err != nil {
		return nil, err
	}
	if input.CategoryID == 0 {
		return nil, ErrCategoryNotFound
	}
	if _, err := s.categoryRepo.GetByID(ctx, input.CategoryID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	stock, err := stockForCreate(input.Stock)
	if err != nil {
		return nil, err
	}

	cid := input.CategoryID
	product := &models.Product{
		Title:       input.Title,
		Description: input.Description,
		BasePrice:   input.Price,
		ImgURL:      input.ImgURL,
		ProductType: "extra",
		Stock:       stock,
		CategoryID:  &cid,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	loaded, err := s.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		return nil, err
	}
	return toExtraResponse(loaded), nil
}

func (s *ProductService) UpdateExtra(ctx context.Context, id uint, input dto.UpdateExtraInput) (*dto.ExtraResponse, error) {
	if input.Title == nil && input.Description == nil && input.Price == nil && input.ImgURL == nil && input.CategoryID == nil && input.Stock == nil {
		return nil, ErrNoFieldsToUpdate
	}

	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrProductNotFound
		}
		return nil, err
	}
	if product.ProductType != "extra" {
		return nil, ErrNotExtraProduct
	}

	if input.CategoryID != nil {
		if *input.CategoryID == 0 {
			return nil, ErrCategoryNotFound
		}
		if _, err := s.categoryRepo.GetByID(ctx, *input.CategoryID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrCategoryNotFound
			}
			return nil, err
		}
		c := *input.CategoryID
		product.CategoryID = &c
	}

	needsProductValidation := input.Title != nil || input.Description != nil || input.Price != nil
	if needsProductValidation {
		title := product.Title
		description := product.Description
		price := product.BasePrice
		if input.Title != nil {
			title = *input.Title
		}
		if input.Description != nil {
			description = *input.Description
		}
		if input.Price != nil {
			price = *input.Price
		}
		if err := validateSimpleProduct(title, description, price); err != nil {
			return nil, err
		}
	}

	if input.Title != nil {
		product.Title = *input.Title
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Price != nil {
		product.BasePrice = *input.Price
	}
	if input.ImgURL != nil {
		product.ImgURL = *input.ImgURL
	}
	if input.Stock != nil {
		if *input.Stock < 0 {
			return nil, ErrInvalidStock
		}
		product.Stock = *input.Stock
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, err
	}

	loaded, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toExtraResponse(loaded), nil
}

func (s *ProductService) DeleteExtra(ctx context.Context, id uint) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrProductNotFound
		}
		return err
	}
	if product.ProductType != "extra" {
		return ErrNotExtraProduct
	}
	return s.productRepo.Delete(ctx, id)
}



func validatePizzaInput(title, description string, price float64) error {
	if len(strings.TrimSpace(title)) < 6 || len(title) > 100 {
		return ErrInvalidTitle
	}
	if len(strings.TrimSpace(description)) < 6 || len(description) > 100 {
		return ErrInvalidDesc
	}
	if price < 100 || price > 10000 {
		return ErrInvalidPrice
	}
	return nil
}

func validateDrinkInput(title, description string, price float64) error {
	if len(strings.TrimSpace(title)) < 3 || len(title) > 50 {
		return ErrInvalidTitle
	}
	if len(strings.TrimSpace(description)) < 3 || len(description) > 100 {
		return ErrInvalidDesc
	}
	if price < 1 || price > 5000 {
		return ErrInvalidPrice
	}
	return nil
}

func validateSimpleProduct(title, description string, price float64) error {
	if len(strings.TrimSpace(title)) < 3 || len(title) > 50 {
		return ErrInvalidTitle
	}
	if len(strings.TrimSpace(description)) < 3 || len(description) > 50 {
		return ErrInvalidDesc
	}
	if price < 1 || price > 5000 {
		return ErrInvalidPrice
	}
	return nil
}

func stockForCreate(p *int) (int, error) {
	if p == nil {
		return 5, nil
	}
	if *p < 0 {
		return 0, ErrInvalidStock
	}
	return *p, nil
}

// Маппинг в DTO

func toPizzaResponse(p *models.Product, pizza *models.Pizza) *dto.PizzaResponse {
	return &dto.PizzaResponse{
		ID:             p.ID,
		Title:          p.Title,
		Description:    p.Description,
		Price:          p.BasePrice,
		ImgURL:         p.ImgURL,
		Stock:          p.Stock,
		SizePizzaID:   pizza.SizePizzaID,
		SizePizza:      pizza.SizePizza.Title,
		ThicknessPizzaID: pizza.ThicknessPizzaID,
		ThicknessPizza: pizza.ThicknessPizza.Title,
	}
}

func toDrinkResponse(p *models.Product, drink *models.Drink) *dto.DrinkResponse {
	return &dto.DrinkResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		Price:       p.BasePrice,
		ImgURL:      p.ImgURL,
		Stock:       p.Stock,
		TypeDrinkID: drink.TypeDrinkID,
		TypeDrink:   drink.TypeDrink.Title,
		VolumeID:    drink.VolumeID,
		Volume:      drink.VolumeDrink.Title,
	}
}

func toExtraResponse(p *models.Product) *dto.ExtraResponse {
	var catID uint
	var catTitle string
	if p.CategoryID != nil {
		catID = *p.CategoryID
	}
	if p.Category != nil {
		catTitle = p.Category.Title
	}
	return &dto.ExtraResponse{
		ID:            p.ID,
		Title:         p.Title,
		Description:   p.Description,
		Price:         p.BasePrice,
		ImgURL:        p.ImgURL,
		Stock:         p.Stock,
		CategoryID:    catID,
		CategoryTitle: catTitle,
	}
}

func (s *ProductService) toProductResponse(ctx context.Context, p *models.Product) (*dto.ProductResponse, error) {
	resp := &dto.ProductResponse{
		ID:          p.ID,
		Title:       p.Title,
		Description: p.Description,
		BasePrice:   p.BasePrice,
		ImgURL:      p.ImgURL,
		Stock:       p.Stock,
		ProductType: p.ProductType,
	}

	// Нормализация для фронтенда (в интерфейсе ожидаются plural-значения).
	switch resp.ProductType {
	case "drink":
		resp.ProductType = "drinks"
	case "pizza":
		resp.ProductType = "pizzas"
	case "extra":
		resp.ProductType = "extras"
		if p.CategoryID != nil {
			resp.CategoryID = *p.CategoryID
		}
		if p.Category != nil {
			resp.CategoryTitle = p.Category.Title
		}
	}

	switch p.ProductType {
	case "drink", "drinks":
		drink, err := s.productRepo.GetDrinkByProductID(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		resp.TypeDrink = drink.TypeDrink.Title
		resp.Volume = drink.VolumeDrink.Title
	case "pizza", "pizzas":
		pizza, err := s.productRepo.GetPizzaByProductID(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		resp.SizePizza = pizza.SizePizza.Title
		resp.ThicknessPizza = pizza.ThicknessPizza.Title
	}

	return resp, nil
}