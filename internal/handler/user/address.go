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

type AddressHandler struct {
	svc *user.AddressService
}

func NewAddressHandler(svc *user.AddressService) *AddressHandler {
	return &AddressHandler{svc: svc}
}

func (h *AddressHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}
	list, err := h.svc.List(r.Context(), userID)
	if err != nil {
		utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, list, http.StatusOK)
}

func (h *AddressHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}
	var in dto.CreateDeliveryAddressInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		utils.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}
	if in.Line == "" {
		utils.JsonError(w, "укажите адрес (поле line)", http.StatusBadRequest)
		return
	}
	out, err := h.svc.Create(r.Context(), userID, in)
	if err != nil {
		utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	utils.JsonResponse(w, out, http.StatusCreated)
}

func (h *AddressHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}
	id, err := parseID(r)
	if err != nil {
		utils.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}
	var in dto.UpdateDeliveryAddressInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		utils.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}
	if in.Label == nil && in.Line == nil && in.IsDefault == nil {
		utils.JsonError(w, "нет полей для обновления", http.StatusBadRequest)
		return
	}
	out, err := h.svc.Update(r.Context(), userID, id, in)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrAddressNotFound):
			utils.JsonError(w, err.Error(), http.StatusNotFound)
		case errors.Is(err, user.ErrAddressNotYours):
			utils.JsonError(w, err.Error(), http.StatusForbidden)
		case errors.Is(err, user.ErrAddressLineEmpty):
			utils.JsonError(w, err.Error(), http.StatusBadRequest)
		default:
			utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}
	utils.JsonResponse(w, out, http.StatusOK)
}

func (h *AddressHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}
	id, err := parseID(r)
	if err != nil {
		utils.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}
	if err := h.svc.Delete(r.Context(), userID, id); err != nil {
		switch {
		case errors.Is(err, user.ErrAddressNotFound):
			utils.JsonError(w, err.Error(), http.StatusNotFound)
		default:
			utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}
	utils.JsonResponse(w, map[string]string{"message": "deleted"}, http.StatusOK)
}
