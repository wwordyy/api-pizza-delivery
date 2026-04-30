package userService

import (
	"api-pizza-delivery/internal/dto"
	"api-pizza-delivery/internal/mail"
	"api-pizza-delivery/internal/models"
	user "api-pizza-delivery/internal/repository/user"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	ErrOrderNotFound = errors.New("заказ не найден")
	ErrEmptyCart     = errors.New("корзина пуста")
	ErrInvalidDate   = errors.New("дата доставки не может быть в прошлом")
	ErrBankCardEmpty = errors.New("при оплате банковской картой поле bank_card обязательно")
)

type OrderService struct {
	orderRepo    *user.OrderRepository
	cartRepo     *user.CartRepository
	profileRepo  *user.ProfileRepository
	mailSender   mail.Sender
	db           *gorm.DB
}

func NewOrderService(
	orderRepo *user.OrderRepository,
	cartRepo *user.CartRepository,
	profileRepo *user.ProfileRepository,
	mailSender mail.Sender,
	db *gorm.DB,
) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		cartRepo:    cartRepo,
		profileRepo: profileRepo,
		mailSender:  mailSender,
		db:          db,
	}
}

func (s *OrderService) GetOrder(ctx context.Context, userID uint, orderID uint) (*dto.OrderResponse, error) {
	o, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		if errors.Is(err, user.ErrOrderNotFound) {
			return nil, ErrOrderNotFound
		}
		return nil, err
	}
	if o.UserID != userID {
		return nil, ErrOrderNotFound
	}
	resp := toOrderResponse(*o)
	return &resp, nil
}

func (s *OrderService) GetOrders(ctx context.Context, userID uint) ([]dto.OrderResponse, error) {
	orders, err := s.orderRepo.GetAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := make([]dto.OrderResponse, len(orders))
	for i, o := range orders {
		response[i] = toOrderResponse(o)
	}
	return response, nil
}

func (s *OrderService) CreateOrder(ctx context.Context, userID uint, input dto.CreateOrderInput) (*dto.OrderResponse, error) {
	if input.DeliveryDate != "" {
		deliveryDate, err := time.Parse("2006-01-02", input.DeliveryDate)
		if err != nil {
			return nil, ErrInvalidDate
		}
		if deliveryDate.Before(time.Now().Truncate(24 * time.Hour)) {
			return nil, ErrInvalidDate
		}
	}

	var paymentMethod models.PaymentMethod
	if err := s.db.WithContext(ctx).First(&paymentMethod, input.PaymentMethodID).Error; err != nil {
		return nil, err
	}
	if paymentMethod.Title == "Банковской картой" && strings.TrimSpace(input.BankCard) == "" {
		return nil, ErrBankCardEmpty
	}

	cart, err := s.cartRepo.GetWithItems(ctx, userID)
	if err != nil || len(cart.CartItems) == 0 {
		return nil, ErrEmptyCart
	}

	receiptLines := make([]mail.ReceiptLine, 0, len(cart.CartItems))
	var total float64
	for _, item := range cart.CartItems {
		sub := item.Product.BasePrice * float64(item.Quantity)
		total += sub
		receiptLines = append(receiptLines, mail.ReceiptLine{
			Title:    item.Product.Title,
			Quantity: item.Quantity,
			Price:    item.Product.BasePrice,
			Subtotal: sub,
		})
	}

	order := &models.Order{
		UserID:          userID,
		DateOfOrder:     time.Now(),
		SumOrder:        total,
		Addresses:       input.Addresses,
		PaymentMethodID: input.PaymentMethodID,
		StatusOrderID:   1,
	}

	if input.DeliveryDate != "" {
		deliveryDate, _ := time.Parse("2006-01-02", input.DeliveryDate)
		order.DeliveryDate = &deliveryDate
	}

	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, item := range cart.CartItems {
			res := tx.Model(&models.Product{}).
				Where("id = ? AND stock >= ?", item.ProductID, item.Quantity).
				Update("stock", gorm.Expr("stock - ?", item.Quantity))
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected != 1 {
				return ErrInsufficientStock
			}
		}

		if err := tx.Create(order).Error; err != nil {
			return err
		}
		for _, item := range cart.CartItems {
			oi := models.OrderItem{
				OrderID:   order.ID,
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.Product.BasePrice,
			}
			if err := tx.Create(&oi).Error; err != nil {
				return err
			}
		}
		return tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error
	})
	if err != nil {
		return nil, err
	}

	created, err := s.orderRepo.GetByID(ctx, order.ID)
	if err != nil {
		return nil, err
	}

	if s.mailSender != nil {
		u, errU := s.profileRepo.GetByID(ctx, userID)
		if errU == nil {
			rd := mail.ReceiptData{
				OrderID:       created.ID,
				Username:      u.Username,
				DateOfOrder:   created.DateOfOrder,
				DeliveryDate:  created.DeliveryDate,
				Address:       created.Addresses,
				PaymentMethod: created.PaymentMethod.Title,
				Status:        created.StatusOrder.Title,
				Lines:         receiptLines,
				Total:         created.SumOrder,
			}
			go func(email string, data mail.ReceiptData) {
				ctxM, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				if err := s.mailSender.SendOrderReceipt(ctxM, email, data); err != nil {
					log.Printf("order receipt email: %v", err)
				}
			}(u.Email, rd)
		}
	}

	resp := toOrderResponse(*created)
	return &resp, nil
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
