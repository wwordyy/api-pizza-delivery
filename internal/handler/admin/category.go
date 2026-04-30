package adminHandler

import (
	admin "api-pizza-delivery/internal/service/admin"
	u "api-pizza-delivery/internal/utils"
	"encoding/json"
	"errors"
	"net/http"

	"api-pizza-delivery/internal/dto"
)

type CategoryHandler struct {
	categoryService *admin.CategoryService
}

func NewCategoryHandler(categoryService *admin.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	list, err := h.categoryService.List(r.Context())
	if err != nil {
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	u.JsonResponse(w, list, http.StatusOK)
}

func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateProductCategoryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}
	resp, err := h.categoryService.Create(r.Context(), input)
	if err != nil {
		handleCategoryError(w, err)
		return
	}
	u.JsonResponse(w, resp, http.StatusCreated)
}

func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		u.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}
	var input dto.UpdateProductCategoryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}
	resp, err := h.categoryService.Update(r.Context(), id, input)
	if err != nil {
		handleCategoryError(w, err)
		return
	}
	u.JsonResponse(w, resp, http.StatusOK)
}

func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		u.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}
	if err := h.categoryService.Delete(r.Context(), id); err != nil {
		handleCategoryError(w, err)
		return
	}
	u.JsonResponse(w, map[string]string{"message": "category deleted successfully"}, http.StatusOK)
}

func handleCategoryError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, admin.ErrCategoryNotFound):
		u.JsonError(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, admin.ErrCategoryInUse):
		u.JsonError(w, err.Error(), http.StatusConflict)
	case errors.Is(err, admin.ErrCategoryTitle):
		u.JsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, admin.ErrNoFieldsToUpdate):
		u.JsonError(w, err.Error(), http.StatusBadRequest)
	default:
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}
