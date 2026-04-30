package userService

import (
	"api-pizza-delivery/internal/dto"
	"api-pizza-delivery/internal/models"
	user "api-pizza-delivery/internal/repository/user"
	"context"
	"errors"
)

var (
	ErrReviewNotFound  = errors.New("отзыв не найден")
	ErrNotOrdered      = errors.New("отзыв можно оставить только на заказанный товар")
	ErrAlreadyReviewed = errors.New("вы уже оставляли отзыв на этот товар")
	ErrInvalidRating   = errors.New("оценка от 1 до 5")
	ErrCommentTooLong  = errors.New("комментарий не длиннее 200 символов")
)

type ReviewService struct {
	reviewRepo *user.ReviewRepository
}

func NewReviewService(reviewRepo *user.ReviewRepository) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo}
}

func (s *ReviewService) GetByProductID(ctx context.Context, productID uint) ([]dto.ReviewResponse, error) {
	reviews, err := s.reviewRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.ReviewResponse, len(reviews))
	for i, r := range reviews {
		response[i] = dto.ReviewResponse{
			ID:        r.ID,
			UserID:    r.UserID,
			ProductID: r.ProductID,
			Rating:    r.Rating,
			Comment:   r.Comment,
			CreatedAt: r.CreatedAt,
		}
	}
	return response, nil
}

func (s *ReviewService) Create(ctx context.Context, userID, productID uint, input dto.CreateReviewInput) (*dto.ReviewResponse, error) {

	if input.Rating < 1 || input.Rating > 5 {
		return nil, ErrInvalidRating
	}

	if len(input.Comment) > 200 {
		return nil, ErrCommentTooLong
	}

	ordered, err := s.reviewRepo.HasOrdered(ctx, userID, productID)
	if err != nil {
		return nil, err
	}
	if !ordered {
		return nil, ErrNotOrdered
	}

	reviewed, err := s.reviewRepo.AlreadyReviewed(ctx, userID, productID)
	if err != nil {
		return nil, err
	}
	if reviewed {
		return nil, ErrAlreadyReviewed
	}

	review := &models.Review{
		UserID:    userID,
		ProductID: productID,
		Rating:    input.Rating,
		Comment:   input.Comment,
	}

	if err := s.reviewRepo.Create(ctx, review); err != nil {
		return nil, err
	}

	return &dto.ReviewResponse{
		ID:        review.ID,
		UserID:    review.UserID,
		ProductID: review.ProductID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt,
	}, nil
}
