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



type CartHandler struct {
	cartService *user.CartService
}

func NewCartHandler(cartService *user.CartService) *CartHandler {
	return &CartHandler{cartService: cartService}
}

// GET /api/
func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	cart, err := h.cartService.GetCart(r.Context(), userID)
	if err != nil {
		utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, cart, http.StatusOK)
}

// POST /api/cart
func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	var input dto.AddToCartInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if input.ProductID == 0 {
		utils.JsonError(w, "укажите product_id", http.StatusBadRequest)
		return
	}
	if input.Quantity <= 0 {
		utils.JsonError(w, "количество должно быть больше 0", http.StatusBadRequest)
		return
	}

	if err := h.cartService.AddItem(r.Context(), userID, input.ProductID, input.Quantity); err != nil {
		handleCartError(w, err)
		return
	}

	utils.JsonResponse(w, map[string]string{"message": "item added to cart"}, http.StatusOK)
}

// PATCH /api/cart/{id}
func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	itemID, err := parseID(r)
	if err != nil {
		utils.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}

	var input dto.UpdateCartItemInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if err := h.cartService.UpdateItem(r.Context(), userID, itemID, input.Quantity); err != nil {
		handleCartError(w, err)
		return
	}

	utils.JsonResponse(w, map[string]string{"message": "cart item updated"}, http.StatusOK)
}

// DELETE /api/cart/{id}
func (h *CartHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	itemID, err := parseID(r)
	if err != nil {
		utils.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}

	if err := h.cartService.DeleteItem(r.Context(), userID, itemID); err != nil {
		handleCartError(w, err)
		return
	}

	utils.JsonResponse(w, map[string]string{"message": "item removed from cart"}, http.StatusOK)
}

func handleCartError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, user.ErrCartItemNotFound):
		utils.JsonError(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, user.ErrInvalidQuantity):
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, user.ErrNotYourItem):
		utils.JsonError(w, err.Error(), http.StatusForbidden)
	case errors.Is(err, user.ErrInsufficientStock):
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, user.ErrProductNotFoundCart):
		utils.JsonError(w, err.Error(), http.StatusNotFound)
	default:
		utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}

func parseID(r *http.Request) (uint, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
