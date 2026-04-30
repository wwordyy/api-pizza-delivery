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
	ErrCategoryNotFound = errors.New("категория не найдена")
	ErrCategoryInUse    = errors.New("категория используется товарами")
	ErrCategoryTitle    = errors.New("неверное название категории")
)

type CategoryService struct {
	categoryRepo *adminRepo.CategoryRepository
}

func NewCategoryService(categoryRepo *adminRepo.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) List(ctx context.Context) ([]dto.ProductCategoryResponse, error) {
	list, err := s.categoryRepo.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]dto.ProductCategoryResponse, 0, len(list))
	for i := range list {
		out = append(out, dto.ProductCategoryResponse{
			ID:    list[i].ID,
			Title: list[i].Title,
		})
	}
	return out, nil
}

func (s *CategoryService) Create(ctx context.Context, input dto.CreateProductCategoryInput) (*dto.ProductCategoryResponse, error) {
	title := strings.TrimSpace(input.Title)
	if len(title) < 2 || len(title) > 80 {
		return nil, ErrCategoryTitle
	}
	c := &models.ProductCategory{Title: title}
	if err := s.categoryRepo.Create(ctx, c); err != nil {
		return nil, err
	}
	return &dto.ProductCategoryResponse{ID: c.ID, Title: c.Title}, nil
}

func (s *CategoryService) Update(ctx context.Context, id uint, input dto.UpdateProductCategoryInput) (*dto.ProductCategoryResponse, error) {
	if input.Title == nil {
		return nil, ErrNoFieldsToUpdate
	}
	c, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}
	t := strings.TrimSpace(*input.Title)
	if len(t) < 2 || len(t) > 80 {
		return nil, ErrCategoryTitle
	}
	c.Title = t
	if err := s.categoryRepo.Update(ctx, c); err != nil {
		return nil, err
	}
	return &dto.ProductCategoryResponse{ID: c.ID, Title: c.Title}, nil
}

func (s *CategoryService) Delete(ctx context.Context, id uint) error {
	_, err := s.categoryRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCategoryNotFound
		}
		return err
	}
	n, err := s.categoryRepo.CountProducts(ctx, id)
	if err != nil {
		return err
	}
	if n > 0 {
		return ErrCategoryInUse
	}
	return s.categoryRepo.Delete(ctx, id)
}
