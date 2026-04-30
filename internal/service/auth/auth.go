package authService

import (
	"api-pizza-delivery/internal/dto"
	"api-pizza-delivery/internal/mail"
	"api-pizza-delivery/internal/models"
	auth "api-pizza-delivery/internal/repository/auth"
	"api-pizza-delivery/internal/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo   *auth.UserRepository
	resetRepo  *auth.PasswordResetRepository
	mailSender mail.Sender
	jwtSecret  string
	jwtExpire  time.Duration
}

var (
	ErrEmailTaken     = errors.New("этот email уже зарегистрирован")
	ErrInvalidCreds   = errors.New("неверный email или пароль")
	ErrInvalidCode    = errors.New("неверный или просроченный код сброса")
	ErrMailSendFailed = errors.New("не удалось отправить письмо")
)

func NewAuthService(
	userRepo *auth.UserRepository,
	resetRepo *auth.PasswordResetRepository,
	mailSender mail.Sender,
	jwtSecret string,
	jwtExpireHours int) *AuthService {

	if jwtExpireHours <= 0 {
		jwtExpireHours = 24
	}

	return &AuthService{
		userRepo:   userRepo,
		resetRepo:  resetRepo,
		mailSender: mailSender,
		jwtSecret:  jwtSecret,
		jwtExpire:  time.Duration(jwtExpireHours) * time.Hour,
	}
}

func (s *AuthService) Register(ctx context.Context, input dto.RegisterInput) (*dto.AuthResponse, error) {

	exists, err := s.userRepo.EmailExists(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("Checking email: %w", err)
	}

	if exists {
		return nil, ErrEmailTaken
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Generating password: %w", err)
	}

	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Phone:    input.Phone,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("Creating user: %w", err)
	}

	user, err = s.userRepo.FindByEmail(ctx, user.Email)
	if err != nil {
		return nil, fmt.Errorf("Finding user by email: %w", err)
	}

	token, err := utils.GenerateToken(user, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{
		Token: token,
		User:  user,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input dto.LoginInput) (*dto.AuthResponse, error) {

	
	user, err := s.userRepo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, ErrInvalidCreds
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, ErrInvalidCreds
	}

	token, err := utils.GenerateToken(user, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, err
	}

	if s.mailSender != nil {
		email, username := user.Email, user.Username
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
			defer cancel()
			if err := s.mailSender.SendLoginNotification(ctx, email, username, time.Now()); err != nil {
				log.Printf("login notification email: %v", err)
			}
		}()
	}

	return &dto.AuthResponse{
			Token: token,
			User:  user,
		},
		nil
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, input dto.ResetRequestInput) (string, error) {

	exists, err := s.userRepo.EmailExists(ctx, input.Email)
	if err != nil {
		return "", fmt.Errorf("checking email: %w", err)
	}
	if !exists {
		return "", nil
	}

	code := utils.GenerateCode()
	expiresAt := time.Now().Add(15 * time.Minute)

	if err := s.resetRepo.Create(ctx, input.Email, code, expiresAt); err != nil {
		return "", fmt.Errorf("creating reset code: %w", err)
	}

	if s.mailSender != nil {
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		if err := s.mailSender.SendPasswordResetCode(ctx, input.Email, code); err != nil {
			log.Printf("password reset: SMTP send failed for %q: %v", input.Email, err)
			return "", fmt.Errorf("%w: %v", ErrMailSendFailed, err)
		}
		return "", nil
	}

	// Без SMTP (локальная разработка): код можно вернуть клиенту для тестов
	return code, nil
}


func (s *AuthService) ResetPassword(ctx context.Context, input dto.ResetPasswordInput) error {
	reset, err := s.resetRepo.FindValid(ctx, input.Code, input.Email)
	
	if err != nil {
		return ErrInvalidCode
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, input.Email, string(hashedPassword)); err != nil {
		return fmt.Errorf("updating password: %w", err)
	}

	if err := s.resetRepo.MarkUsed(ctx, reset.ID); err != nil {
		return fmt.Errorf("marking code as used: %w", err)
	}

	return nil
}