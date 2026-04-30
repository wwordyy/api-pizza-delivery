package mail

import (
	"context"
	"encoding/base64"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
}

func (c SMTPConfig) IsConfigured() bool {
	return c.Host != "" && c.Port != "" && c.User != "" && c.Password != ""
}

type SMTPSender struct {
	cfg SMTPConfig
}

func NewSMTP(cfg SMTPConfig) *SMTPSender {
	if cfg.From == "" {
		cfg.From = cfg.User
	}
	return &SMTPSender{cfg: cfg}
}

func encodeSubjectUTF8(s string) string {
	return fmt.Sprintf("=?UTF-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(s)))
}

func (s *SMTPSender) sendTextEmail(ctx context.Context, to, subjectEncoded, body string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("Date: " + time.Now().Format(time.RFC1123Z) + "\r\n")
	fmt.Fprintf(&b, "From: %s\r\n", s.cfg.From)
	fmt.Fprintf(&b, "To: %s\r\n", to)
	fmt.Fprintf(&b, "Subject: %s\r\n", subjectEncoded)
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("\r\n")
	b.WriteString(body)

	addr := net.JoinHostPort(s.cfg.Host, s.cfg.Port)
	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)
	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, []byte(b.String()))
}

func (s *SMTPSender) SendPasswordResetCode(ctx context.Context, to, code string) error {
	subj := encodeSubjectUTF8("Код восстановления пароля")
	body := fmt.Sprintf(
		"Здравствуйте!\n\n"+
			"Ваш код для сброса пароля: %s\n\n"+
			"Код действителен 15 минут.\n"+
			"Если вы не запрашивали сброс пароля, проигнорируйте это письмо.\n\n"+
			"С уважением,\nСлужба доставки пиццы",
		code,
	)
	if err := s.sendTextEmail(ctx, to, subj, body); err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}
	return nil
}

func (s *SMTPSender) SendLoginNotification(ctx context.Context, to, username string, at time.Time) error {
	subj := encodeSubjectUTF8("Вход в аккаунт")
	body := fmt.Sprintf(
		"Здравствуйте, %s!\n\n"+
			"В ваш аккаунт только что выполнен вход.\n"+
			"Время: %s\n\n"+
			"Если это были не вы, срочно смените пароль в настройках профиля.\n\n"+
			"С уважением,\nСлужба доставки пиццы",
		username,
		at.Format("02.01.2006 15:04:05 MST"),
	)
	if err := s.sendTextEmail(ctx, to, subj, body); err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}
	return nil
}

func (s *SMTPSender) SendOrderReceipt(ctx context.Context, to string, data ReceiptData) error {
	subj := encodeSubjectUTF8(fmt.Sprintf("Чек заказа №%d", data.OrderID))
	var b strings.Builder
	fmt.Fprintf(&b, "Здравствуйте, %s!\n\n", data.Username)
	fmt.Fprintf(&b, "Спасибо за заказ! Ниже данные электронного чека.\n\n")
	fmt.Fprintf(&b, "Номер заказа: %d\n", data.OrderID)
	fmt.Fprintf(&b, "Дата оформления: %s\n", data.DateOfOrder.Format("02.01.2006 15:04:05"))
	if data.DeliveryDate != nil {
		fmt.Fprintf(&b, "\nПланируемая дата доставки: %s\n", data.DeliveryDate.Format("02.01.2006"))
	}
	fmt.Fprintf(&b, "\nАдрес доставки: %s\n", data.Address)
	fmt.Fprintf(&b, "Способ оплаты: %s\n", data.PaymentMethod)
	fmt.Fprintf(&b, "Статус: %s\n\n", data.Status)
	b.WriteString("Состав заказа:\n")
	b.WriteString(strings.Repeat("-", 40) + "\n")
	for _, line := range data.Lines {
		fmt.Fprintf(&b, "%s\n", line.Title)
		fmt.Fprintf(&b, "  %d × %.2f ₽ = %.2f ₽\n", line.Quantity, line.Price, line.Subtotal)
	}
	b.WriteString(strings.Repeat("-", 40) + "\n")
	fmt.Fprintf(&b, "Итого: %.2f ₽\n\n", data.Total)
	b.WriteString("С уважением,\nСлужба доставки пиццы\n")
	if err := s.sendTextEmail(ctx, to, subj, b.String()); err != nil {
		return fmt.Errorf("smtp send: %w", err)
	}
	return nil
}
