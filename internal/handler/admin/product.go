package adminHandler

import (
	admin "api-pizza-delivery/internal/service/admin"
	u "api-pizza-delivery/internal/utils"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"api-pizza-delivery/internal/dto"
)

type ProductHandler struct {
	productService *admin.ProductService
}

func NewProductHandler(productService *admin.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}



func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	filter := dto.ProductFilter{
		Type:   r.URL.Query().Get("type"),
		Search: r.URL.Query().Get("search"),
		Sort:   r.URL.Query().Get("sort"),
	}
	if s := r.URL.Query().Get("category_id"); s != "" {
		id, err := strconv.ParseUint(s, 10, 32)
		if err != nil {
			u.JsonError(w, "неверный category_id", http.StatusBadRequest)
			return
		}
		filter.CategoryID = uint(id)
	}

	products, err := h.productService.GetAll(r.Context(), filter)
	if err != nil {
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	u.JsonResponse(w, products, http.StatusOK)
}


func (h *ProductHandler) GetAllPizzas(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetAllPizzas(r.Context())
	if err != nil {
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	u.JsonResponse(w, products, http.StatusOK)
}

func (h *ProductHandler) GetAllDrinks(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetAllDrinks(r.Context())
	if err != nil {
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	u.JsonResponse(w, products, http.StatusOK)
}

func (h *ProductHandler) GetAllExtras(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetAllExtras(r.Context())
	if err != nil {
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	u.JsonResponse(w, products, http.StatusOK)
}

func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		u.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}

	product, err := h.productService.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, admin.ErrProductNotFound) {
			u.JsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	u.JsonResponse(w, product, http.StatusOK)
}

func (h *ProductHandler) CreatePizza(w http.ResponseWriter, r *http.Request) {
	var input dto.CreatePizzaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	resp, err := h.productService.CreatePizza(r.Context(), input)
	if err != nil {
		handleProductError(w, err)
		return
	}

	u.JsonResponse(w, resp, http.StatusCreated)
}

func (h *ProductHandler) UpdatePizza(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		u.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}

	var input dto.UpdatePizzaInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	resp, err := h.productService.UpdatePizza(r.Context(), id, input)
	if err != nil {
		handleProductError(w, err)
		return
	}

	u.JsonResponse(w, resp, http.StatusOK)
}

func (h *ProductHandler) DeletePizza(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		u.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}

	if err := h.productService.DeletePizza(r.Context(), id); err != nil {
		handleProductError(w, err)
		return
	}

	u.JsonResponse(w, map[string]string{"message": "pizza deleted successfully"}, http.StatusOK)
}

func (h *ProductHandler) CreateDrink(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateDrinkInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	resp, err := h.productService.CreateDrink(r.Context(), input)
	if err != nil {
		handleProductError(w, err)
		return
	}

	u.JsonResponse(w, resp, http.StatusCreated)
}

func (h *ProductHandler) UpdateDrink(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		u.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}

	var input dto.UpdateDrinkInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	resp, err := h.productService.UpdateDrink(r.Context(), id, input)
	if err != nil {
		handleProductError(w, err)
		return
	}

	u.JsonResponse(w, resp, http.StatusOK)
}

func (h *ProductHandler) DeleteDrink(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		u.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}

	if err := h.productService.DeleteDrink(r.Context(), id); err != nil {
		handleProductError(w, err)
		return
	}

	u.JsonResponse(w, map[string]string{"message": "drink deleted successfully"}, http.StatusOK)
}

func (h *ProductHandler) CreateExtra(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateExtraInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	resp, err := h.productService.CreateExtra(r.Context(), input)
	if err != nil {
		handleProductError(w, err)
		return
	}

	u.JsonResponse(w, resp, http.StatusCreated)
}

func (h *ProductHandler) UpdateExtra(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		u.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}

	var input dto.UpdateExtraInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	resp, err := h.productService.UpdateExtra(r.Context(), id, input)
	if err != nil {
		handleProductError(w, err)
		return
	}

	u.JsonResponse(w, resp, http.StatusOK)
}

func (h *ProductHandler) DeleteExtra(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		u.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}

	if err := h.productService.DeleteExtra(r.Context(), id); err != nil {
		handleProductError(w, err)
		return
	}

	u.JsonResponse(w, map[string]string{"message": "extra deleted successfully"}, http.StatusOK)
}

func parseID(r *http.Request) (uint, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

func handleProductError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, admin.ErrProductNotFound):
		u.JsonError(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, admin.ErrNotExtraProduct):
		u.JsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, admin.ErrCategoryNotFound):
		u.JsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, admin.ErrNoFieldsToUpdate):
		u.JsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, admin.ErrInvalidPrice):
		u.JsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, admin.ErrInvalidTitle):
		u.JsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, admin.ErrInvalidDesc):
		u.JsonError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, admin.ErrInvalidStock):
		u.JsonError(w, err.Error(), http.StatusBadRequest)
	default:
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
	}
}
