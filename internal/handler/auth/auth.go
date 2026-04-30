package authHandler

import (
	"api-pizza-delivery/internal/dto"
	u "api-pizza-delivery/internal/utils"
	"api-pizza-delivery/internal/middleware"
	auth "api-pizza-delivery/internal/service/auth"
	"encoding/json"
	"errors"
	"net/http"
	
	
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}



func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input dto.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.Email == "" || input.Password == "" {
		u.JsonError(w, "укажите имя пользователя, email и пароль", http.StatusBadRequest)
		return
	}
	if len(input.Password) < 6 {
		u.JsonError(w, "пароль не короче 6 символов", http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Register(r.Context(), input)
	if err != nil {
		if errors.Is(err, auth.ErrEmailTaken) {
			u.JsonError(w, err.Error(), http.StatusConflict)
			return
		}
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	u.JsonResponse(w, resp, http.StatusCreated)
}


func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Password == "" {
		u.JsonError(w, "укажите email и пароль", http.StatusBadRequest)
		return
	}

	resp, err := h.authService.Login(r.Context(), input)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCreds) {
			u.JsonError(w, err.Error(), http.StatusUnauthorized)
			return
		}
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	u.JsonResponse(w, resp, http.StatusOK)
}


func (h *AuthHandler) ResetRequest(w http.ResponseWriter, r *http.Request) {
	var input dto.ResetRequestInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if input.Email == "" {
		u.JsonError(w, "укажите email", http.StatusBadRequest)
		return
	}

	code, err := h.authService.RequestPasswordReset(r.Context(), input)
	if err != nil {
		if errors.Is(err, auth.ErrMailSendFailed) {
			u.JsonError(w, err.Error(), http.StatusBadGateway)
			return
		}
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	resp := map[string]string{"message": "If this email is registered, a reset code has been sent."}
	if code != "" {
		resp["code"] = code
		resp["message"] = "reset code (dev mode without SMTP): use code or configure SMTP_USER/SMTP_PASSWORD"
	}
	u.JsonResponse(w, resp, http.StatusOK)
}


func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var input dto.ResetPasswordInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		u.JsonError(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	if input.Email == "" || input.Code == "" || input.NewPassword == "" {
		u.JsonError(w, "укажите email, код и новый пароль (new_password)", http.StatusBadRequest)
		return
	}
	if len(input.NewPassword) < 6 {
		u.JsonError(w, "пароль не короче 6 символов", http.StatusBadRequest)
		return
	}

	if err := h.authService.ResetPassword(r.Context(), input); err != nil {
		if errors.Is(err, auth.ErrInvalidCode) {
			u.JsonError(w, err.Error(), http.StatusBadRequest)
			return
		}
		u.JsonError(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	u.JsonResponse(w, map[string]string{"message": "password updated successfully"}, http.StatusOK)
}


func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		u.JsonError(w, "требуется авторизация", http.StatusUnauthorized)
		return
	}

	u.JsonResponse(w, map[string]any{"user_id": userID}, http.StatusOK)
}