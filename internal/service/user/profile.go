package userService

import (
	"api-pizza-delivery/internal/dto"
	"api-pizza-delivery/internal/models"
	user "api-pizza-delivery/internal/repository/user"
	"context"
	"errors"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUsername = errors.New("имя пользователя: от 3 до 20 символов")
	ErrInvalidPhone    = errors.New("телефон укажите в формате +7XXXXXXXXXX")
	ErrWrongPassword   = errors.New("неверный текущий пароль")
	ErrInvalidPassword = errors.New("пароль не менее 6 символов, с заглавной буквой и цифрой")
)

type ProfileService struct {
	profileRepo *user.ProfileRepository
}

func NewProfileService(profileRepo *user.ProfileRepository) *ProfileService {
	return &ProfileService{profileRepo: profileRepo}
}

func (s *ProfileService) GetProfile(ctx context.Context, id uint) (*dto.ProfileResponse, error) {
	user, err := s.profileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toProfileResponse(user), nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, id uint, input dto.UpdateProfileInput) (*dto.ProfileResponse, error) {
	if err := validateUsername(input.Username); err != nil {
		return nil, err
	}

	if input.Phone != "" {
		if err := validatePhone(input.Phone); err != nil {
			return nil, err
		}
	}

	user, err := s.profileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user.Username = input.Username
	user.Phone = input.Phone

	if err := s.profileRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	updated, err := s.profileRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return toProfileResponse(updated), nil
}

func (s *ProfileService) ChangePassword(ctx context.Context, id uint, input dto.ChangePasswordInput) error {
	oldPassword, err := s.profileRepo.GetPasswordByID(ctx, id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(oldPassword), []byte(input.OldPassword)); err != nil {
		return ErrWrongPassword
	}

	if err := validatePassword(input.NewPassword); err != nil {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.profileRepo.UpdatePassword(ctx, id, string(hashedPassword))
}

func validateUsername(username string) error {
	username = strings.TrimSpace(username)
	if len(username) < 3 || len(username) > 20 {
		return ErrInvalidUsername
	}
	return nil
}

func validatePhone(phone string) error {
	re := regexp.MustCompile(`^\+7\d{10}$`)
	if !re.MatchString(phone) {
		return ErrInvalidPhone
	}
	return nil
}

func validatePassword(password string) error {
	if len(password) < 6 {
		return ErrInvalidPassword
	}
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
	hasDigit := regexp.MustCompile(`[0-9]`).MatchString(password)
	if !hasUpper || !hasDigit {
		return ErrInvalidPassword
	}
	return nil
}

func toProfileResponse(u *models.User) *dto.ProfileResponse {
	return &dto.ProfileResponse{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Phone:    u.Phone,
		Role:     u.Role.Title,
	}
}
