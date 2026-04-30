package userHandler

import (
	"api-pizza-delivery/internal/dto"
	utils "api-pizza-delivery/internal/utils"
	"api-pizza-delivery/internal/middleware"
	user "api-pizza-delivery/internal/service/user"
	"encoding/json"
	"errors"
	"net/http"
)

type ProfileHandler struct {
	profileService *user.ProfileService
}

func NewProfileHandler(profileService *user.ProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	profile, err := h.profileService.GetProfile(r.Context(), userID)
	if err != nil {
		utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, profile, http.StatusOK)
}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	var input dto.UpdateProfileInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	profile, err := h.profileService.UpdateProfile(r.Context(), userID, input)
	if err != nil {
		handleProfileError(w, err)
		return
	}

	utils.JsonResponse(w, profile, http.StatusOK)
}

func (h *ProfileHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	var input dto.ChangePasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if input.OldPassword == "" || input.NewPassword == "" {
		utils.JsonError(w, "укажите old_password и new_password", http.StatusBadRequest)
		return
	}

	if err := h.profileService.ChangePassword(r.Context(), userID, input); err != nil {
		handleProfileError(w, err)
		return
	}

	utils.JsonResponse(w, map[string]string{"message": "password changed successfully"}, http.StatusOK)
}

func (h *ProfileHandler) Logout(w http.ResponseWriter, r *http.Request) {
	utils.JsonResponse(w, map[string]string{"message": "logged out successfully"}, http.StatusOK)
}

func handleProfileError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, user.ErrInvalidUsername):
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, user.ErrInvalidPhone):
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, user.ErrWrongPassword):
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, user.ErrInvalidPassword):
		utils.JsonError(w, err.Error(), http.StatusBadRequest)
	default:
		utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}