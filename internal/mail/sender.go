package mail

import (
	"context"
	"time"
)

type ReceiptLine struct {
	Title    string
	Quantity int
	Price    float64
	Subtotal float64
}

type ReceiptData struct {
	OrderID        uint
	Username       string
	DateOfOrder    time.Time
	DeliveryDate   *time.Time
	Address        string
	PaymentMethod  string
	Status         string
	Lines          []ReceiptLine
	Total          float64
}

type Sender interface {
	SendPasswordResetCode(ctx context.Context, to, code string) error
	SendLoginNotification(ctx context.Context, to, username string, at time.Time) error
	SendOrderReceipt(ctx context.Context, to string, data ReceiptData) error
}
