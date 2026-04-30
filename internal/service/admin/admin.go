package adminService

import (
	"api-pizza-delivery/internal/dto"
	"api-pizza-delivery/internal/models"
	admin "api-pizza-delivery/internal/repository/admin"
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// utf8BOM — маркер UTF-8 для Excel и других программ на Windows, чтобы кириллица отображалась верно.
var utf8BOM = []byte{0xEF, 0xBB, 0xBF}

func csvWithUTF8BOM(buf *bytes.Buffer) []byte {
	return append(utf8BOM, buf.Bytes()...)
}

var (
	ErrUserNotFound   = errors.New("пользователь не найден")
	ErrReviewNotFound = errors.New("отзыв не найден")
	ErrInvalidRating  = errors.New("оценка от 1 до 5")
	ErrInvalidComment = errors.New("комментарий не длиннее 200 символов")
)

type AdminService struct {
	adminUserRepo *admin.AdminUserRepository
	reviewRepo    *admin.ReviewRepository
	analyticsRepo *admin.AnalyticsRepository
}

func NewAdminService(
	adminUserRepo *admin.AdminUserRepository,
	reviewRepo *admin.ReviewRepository,
	analyticsRepo *admin.AnalyticsRepository,
) *AdminService {
	return &AdminService{
		adminUserRepo: adminUserRepo,
		reviewRepo:    reviewRepo,
		analyticsRepo: analyticsRepo,
	}
}


func (s *AdminService) GetAllUsers(ctx context.Context) ([]dto.UserResponse, error) {
	users, err := s.adminUserRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.UserResponse, len(users))
	for i, u := range users {
		response[i] = toUserResponse(u)
	}
	return response, nil
}


func (s *AdminService) ExportUsersCSV(ctx context.Context) ([]byte, error) {
	users, err := s.adminUserRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Заголовки
	writer.Write([]string{"ID", "Никнейм", "Email", "Телефон", "Роль", "Дата регистрации"})

	// Данные
	for _, u := range users {
		writer.Write([]string{
			fmt.Sprintf("%d", u.ID),
			u.Username,
			u.Email,
			u.Phone,
			u.Role.Title,
			u.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return csvWithUTF8BOM(&buf), nil
}


func (s *AdminService) GetAllReviews(ctx context.Context) ([]dto.AdminReviewResponse, error) {
	reviews, err := s.reviewRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	response := make([]dto.AdminReviewResponse, len(reviews))
	for i, r := range reviews {
		response[i] = toAdminReviewResponse(r)
	}
	return response, nil
}


func (s *AdminService) DeleteReview(ctx context.Context, id uint) error {
	_, err := s.reviewRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrReviewNotFound
		}
		return err
	}
	return s.reviewRepo.Delete(ctx, id)
}


func (s *AdminService) GetAnalytics(ctx context.Context, filter dto.AnalyticsFilter) (*dto.AnalyticsResponse, error) {
	repoFilter, err := toRepoFilter(filter)
	if err != nil {
		return nil, err
	}

	summary, err := s.analyticsRepo.GetSummary(ctx, repoFilter)
	if err != nil {
		return nil, err
	}

	orders, err := s.analyticsRepo.GetOrders(ctx, repoFilter)
	if err != nil {
		return nil, err
	}

	orderResponses := make([]dto.OrderResponse, len(orders))
	for i, o := range orders {
		orderResponses[i] = toOrderResponse(o)
	}

	return &dto.AnalyticsResponse{
		Summary: dto.AnalyticsSummaryDTO{
			TotalSum:     summary.TotalSum,
			TotalCount:   summary.TotalCount,
			AverageOrder: summary.AverageOrder,
		},
		Orders: orderResponses,
	}, nil
}


func (s *AdminService) ExportAnalyticsCSV(ctx context.Context, filter dto.AnalyticsFilter) ([]byte, error) {
	repoFilter, err := toRepoFilter(filter)
	if err != nil {
		return nil, err
	}

	orders, err := s.analyticsRepo.GetOrders(ctx, repoFilter)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Заголовки
	writer.Write([]string{
		"ID", "ID пользователя", "Дата заказа", "Дата доставки",
		"Сумма", "Адрес", "Способ оплаты", "Статус",
	})

	// Данные
	for _, o := range orders {
		deliveryDate := ""
		if o.DeliveryDate != nil {
			deliveryDate = o.DeliveryDate.Format("2006-01-02")
		}
		writer.Write([]string{
			fmt.Sprintf("%d", o.ID),
			fmt.Sprintf("%d", o.UserID),
			o.DateOfOrder.Format("2006-01-02"),
			deliveryDate,
			fmt.Sprintf("%.2f", o.SumOrder),
			o.Addresses,
			o.PaymentMethod.Title,
			o.StatusOrder.Title,
		})
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return csvWithUTF8BOM(&buf), nil
}




func toRepoFilter(filter dto.AnalyticsFilter) (admin.AnalyticsFilter, error) {
	repoFilter := admin.AnalyticsFilter{
		Sort: filter.Sort,
	}

	if filter.From != "" {
		t, err := time.Parse("2006-01-02", filter.From)
		if err != nil {
			return repoFilter, fmt.Errorf("неверный формат даты «с», используйте ГГГГ-ММ-ДД")
		}
		repoFilter.From = &t
	}

	if filter.To != "" {
		t, err := time.Parse("2006-01-02", filter.To)
		if err != nil {
			return repoFilter, fmt.Errorf("неверный формат даты «по», используйте ГГГГ-ММ-ДД")
		}
		repoFilter.To = &t
	}

	if filter.UserID != 0 {
		repoFilter.UserID = &filter.UserID
	}

	if filter.PaymentMethodID != 0 {
		repoFilter.PaymentMethodID = &filter.PaymentMethodID
	}

	return repoFilter, nil
}

func toUserResponse(u models.User) dto.UserResponse {
	return dto.UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Phone:     u.Phone,
		Role:      u.Role.Title,
		CreatedAt: u.CreatedAt,
	}
}

func toAdminReviewResponse(r models.Review) dto.AdminReviewResponse {
	return dto.AdminReviewResponse{
		ID:          r.ID,
		UserID:      r.UserID,
		ProductID:   r.ProductID,
		ProductName: r.Product.Title,
		Rating:      r.Rating,
		Comment:     r.Comment,
		CreatedAt:   r.CreatedAt,
	}
}


func toOrderResponse(o models.Order) dto.OrderResponse {
	resp := dto.OrderResponse{
		ID:            o.ID,
		UserID:        o.UserID,
		DateOfOrder:   o.DateOfOrder,
		SumOrder:      o.SumOrder,
		Addresses:     o.Addresses,
		PaymentMethod: o.PaymentMethod.Title,
		StatusOrder:   o.StatusOrder.Title,
		CreatedAt:     o.CreatedAt,
	}
	if o.DeliveryDate != nil {
		t := *o.DeliveryDate
		resp.DeliveryDate = &t
	}
	return resp
}
