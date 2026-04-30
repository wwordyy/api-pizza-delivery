package userService

import (
	"api-pizza-delivery/internal/dto"
	user "api-pizza-delivery/internal/repository/user"
	"context"
	"errors"
)

var (
	ErrFavoriteNotFound = errors.New("избранное не найдено")
	ErrAlreadyFavorited = errors.New("товар уже в избранном")
	ErrNotYourFavorite  = errors.New("эта запись избранного вам не принадлежит")
)

type FavoriteService struct {
	favoriteRepo *user.FavoriteRepository
}

func NewFavoriteService(favoriteRepo *user.FavoriteRepository) *FavoriteService {
	return &FavoriteService{favoriteRepo: favoriteRepo}
}

func (s *FavoriteService) GetAll(ctx context.Context, userID uint) ([]dto.FavoriteResponse, error) {
	favorites, err := s.favoriteRepo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.FavoriteResponse, len(favorites))
	for i, f := range favorites {
		response[i] = dto.FavoriteResponse{
			ID:          f.ID,
			ProductID:   f.ProductID,
			ProductName: f.Product.Title,
			Price:       f.Product.BasePrice,
			ImgURL:      f.Product.ImgURL,
		}
	}
	return response, nil
}

func (s *FavoriteService) Add(ctx context.Context, userID, productID uint) error {
	exists, err := s.favoriteRepo.Exists(ctx, userID, productID)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyFavorited
	}
	return s.favoriteRepo.Add(ctx, userID, productID)
}

func (s *FavoriteService) Delete(ctx context.Context, userID, favoriteID uint) error {
	favorite, err := s.favoriteRepo.GetByID(ctx, favoriteID)
	if err != nil {
		return ErrFavoriteNotFound
	}

	if favorite.UserID != userID {
		return ErrNotYourFavorite
	}

	return s.favoriteRepo.Delete(ctx, favoriteID)
}
