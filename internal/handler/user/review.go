package userHandler

import (
	"api-pizza-delivery/internal/dto"
	"api-pizza-delivery/internal/middleware"
	user "api-pizza-delivery/internal/service/user"
	"api-pizza-delivery/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)



type ReviewHandler struct {
	reviewService *user.ReviewService
}

func NewReviewHandler(reviewService *user.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

func (h *ReviewHandler) GetByProductID(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		utils.JsonError(w, "неверный идентификатор товара", http.StatusBadRequest)
		return
	}

	reviews, err := h.reviewService.GetByProductID(r.Context(), uint(productID))
	if err != nil {
		utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, reviews, http.StatusOK)
}

func (h *ReviewHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	productIDStr := chi.URLParam(r, "id")
	productID, err := strconv.ParseUint(productIDStr, 10, 32)
	if err != nil {
		utils.JsonError(w, "неверный идентификатор товара", http.StatusBadRequest)
		return
	}

	var input dto.CreateReviewInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	review, err := h.reviewService.Create(r.Context(), userID, uint(productID), input)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrInvalidRating):
			utils.JsonError(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, user.ErrCommentTooLong):
			utils.JsonError(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, user.ErrNotOrdered):
			utils.JsonError(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, user.ErrAlreadyReviewed):
			utils.JsonError(w, err.Error(), http.StatusConflict)
		default:
			utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	utils.JsonResponse(w, review, http.StatusCreated)
}
