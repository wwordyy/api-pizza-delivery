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


type OrderHandler struct {
	orderService *user.OrderService
}

func NewOrderHandler(orderService *user.OrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	orders, err := h.orderService.GetOrders(r.Context(), userID)
	if err != nil {
		utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	utils.JsonResponse(w, orders, http.StatusOK)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}
	orderID, err := parseID(r)
	if err != nil {
		utils.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}
	order, err := h.orderService.GetOrder(r.Context(), userID, orderID)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrOrderNotFound):
			utils.JsonError(w, err.Error(), http.StatusNotFound)
		default:
			utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}
	utils.JsonResponse(w, order, http.StatusOK)
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		utils.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	var input dto.CreateOrderInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		utils.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if input.Addresses == "" {
		utils.JsonError(w, "укажите адрес доставки (addresses)", http.StatusBadRequest)
		return
	}
	if input.PaymentMethodID == 0 {
		utils.JsonError(w, "укажите payment_method_id", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.CreateOrder(r.Context(), userID, input)
	if err != nil {
		switch {
		case errors.Is(err, user.ErrEmptyCart):
			utils.JsonError(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, user.ErrInvalidDate):
			utils.JsonError(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, user.ErrInsufficientStock):
			utils.JsonError(w, err.Error(), http.StatusBadRequest)
		case errors.Is(err, user.ErrBankCardEmpty):
			utils.JsonError(w, err.Error(), http.StatusBadRequest)
		default:
			utils.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		}
		return
	}

	utils.JsonResponse(w, order, http.StatusCreated)
}
