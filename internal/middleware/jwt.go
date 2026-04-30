package middleware

import (
	"api-pizza-delivery/internal/utils"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const (
	ContextUserID   contextKey = "userID"
	ContextUserRole contextKey = "userRole"
)

func Auth(jwtSecret string) func(http.Handler) http.Handler {
	
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"отсутствует заголовок Authorization"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error":"неверный формат авторизации, ожидается Bearer-токен"}`, http.StatusUnauthorized)
				return
			}

			claims, err := utils.ParseToken(parts[1], jwtSecret)
			if err != nil {
				http.Error(w, `{"error":"недействительный или просроченный токен"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), ContextUserID, claims.UserID)
			ctx = context.WithValue(ctx, ContextUserRole, claims.Role)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}


func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(ContextUserRole).(string)
		if !ok || role != "admin" {
			http.Error(w, `{"error":"доступ запрещён"}`, http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}


func GetUserID(r *http.Request) (uint, bool) {
	id, ok := r.Context().Value(ContextUserID).(uint)
	return id, ok
}

func GetUserRole(r *http.Request) (string, bool) {
	role, ok := r.Context().Value(ContextUserRole).(string)
	return role, ok
}