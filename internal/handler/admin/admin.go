package adminHandler

import (
	"api-pizza-delivery/internal/dto"
	admin "api-pizza-delivery/internal/service/admin"
	u "api-pizza-delivery/internal/utils"
	"errors"
	"net/http"
	"strconv"
)

type AdminHandler struct {
	adminService *admin.AdminService
}

func NewAdminHandler(adminService *admin.AdminService) *AdminHandler {
	return &AdminHandler{adminService: adminService}
}


func (h *AdminHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.adminService.GetAllUsers(r.Context())
	if err != nil {
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	u.JsonResponse(w, users, http.StatusOK)
}


func (h *AdminHandler) ExportUsers(w http.ResponseWriter, r *http.Request) {
	data, err := h.adminService.ExportUsersCSV(r.Context())
	if err != nil {
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=users.csv")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}


func (h *AdminHandler) GetAllReviews(w http.ResponseWriter, r *http.Request) {
	reviews, err := h.adminService.GetAllReviews(r.Context())
	if err != nil {
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	u.JsonResponse(w, reviews, http.StatusOK)
}


func (h *AdminHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		u.JsonError(w, "неверный идентификатор", http.StatusBadRequest)
		return
	}

	if err := h.adminService.DeleteReview(r.Context(), id); err != nil {
		if errors.Is(err, admin.ErrReviewNotFound) {
			u.JsonError(w, err.Error(), http.StatusNotFound)
			return
		}
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	u.JsonResponse(w, map[string]string{"message": "review deleted successfully"}, http.StatusOK)
}


func (h *AdminHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	filter := parseAnalyticsFilter(r)

	analytics, err := h.adminService.GetAnalytics(r.Context(), filter)
	if err != nil {
		u.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	u.JsonResponse(w, analytics, http.StatusOK)
}


func (h *AdminHandler) ExportAnalytics(w http.ResponseWriter, r *http.Request) {
	filter := parseAnalyticsFilter(r)

	data, err := h.adminService.ExportAnalyticsCSV(r.Context(), filter)
	if err != nil {
		u.JsonError(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", "attachment; filename=analytics.csv")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}


func parseAnalyticsFilter(r *http.Request) dto.AnalyticsFilter {
	filter := dto.AnalyticsFilter{
		From: r.URL.Query().Get("from"),
		To:   r.URL.Query().Get("to"),
		Sort: r.URL.Query().Get("sort"),
	}

	if userIDStr := r.URL.Query().Get("user_id"); userIDStr != "" {
		if id, err := strconv.ParseUint(userIDStr, 10, 32); err == nil {
			filter.UserID = uint(id)
		}
	}

	if pmIDStr := r.URL.Query().Get("payment_method_id"); pmIDStr != "" {
		if id, err := strconv.ParseUint(pmIDStr, 10, 32); err == nil {
			filter.PaymentMethodID = uint(id)
		}
	}

	return filter
}
