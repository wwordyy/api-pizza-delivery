package userHandler

import (
	"api-pizza-delivery/internal/dto"
	"api-pizza-delivery/internal/middleware"
	user "api-pizza-delivery/internal/service/user"
	"api-pizza-delivery/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
)



type FavoriteHandler struct {
	favoriteService *user.FavoriteService
}

func NewFavoriteHandler(favoriteService *user.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{favoriteService: favoriteService}
}

func (h *FavoriteHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	favorites, err := h.favoriteService.GetAll(r.Context(), userID)
	if err != nil {
		utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, favorites, http.StatusOK)
}

// POST /api/favorites
func (h *FavoriteHandler) Add(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	var input dto.AddFavoriteInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if input.ProductID == 0 {
		utils.JsonError(w, "укажите product_id", http.StatusBadRequest)
		return
	}

	if err := h.favoriteService.Add(r.Context(), userID, input.ProductID); err != nil {
		switch {
		case errors.Is(err, user.ErrAlreadyFavorited):
			utils.JsonError(w, err.Error(), http.StatusConflict)
		default:
			utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	utils.JsonResponse(w, map[string]string{"message": "added to favorites"}, http.StatusOK)
}

func (h *FavoriteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	favoriteID, err := parseID(r)
	if err != nil {
		utils.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}

	if err := h.favoriteService.Delete(r.Context(), userID, favoriteID); err != nil {
		switch {
		case errors.Is(err, user.ErrFavoriteNotFound):
			utils.JsonError(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, user.ErrNotYourFavorite):
			utils.JsonError(w, err.Error(), http.StatusForbidden)
		default:
			utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	utils.JsonResponse(w, map[string]string{"message": "removed from favorites"}, http.StatusOK)
}
