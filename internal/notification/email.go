package notification

import (
	"fmt"
	"github.com/wneessen/go-mail"
	"os"
)

type (
	emailNotifier struct {
		smtpUser string
		smtpPass string
	}
)

func NewEmailNotifier() Notifier {
	return &emailNotifier{
		smtpUser: os.Getenv("SMTP_USER"),
		smtpPass: os.Getenv("SMTP_PASSWORD"),
	}
}

func (e *emailNotifier) Notify(msg Message, recipient Recipient) error {
	message := mail.NewMsg()
	if err := message.From(e.smtpUser); err != nil {
		return fmt.Errorf("failed to set FROM address: %w", err)
	}
	if err := message.To(recipient.Email); err != nil {
		return fmt.Errorf("failed to set TO address: %w", err)
	}
	message.Subject(msg.Subject)
	message.SetBodyString(mail.TypeTextPlain, msg.Body)

	client, err := mail.NewClient("smtp.gmail.com",
		mail.WithPort(587),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(e.smtpUser),
		mail.WithPassword(e.smtpPass),
	)
	if err != nil {
		return fmt.Errorf("failed to create new mail delivery client: %w", err)
	}

	if err := client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}
