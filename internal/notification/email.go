package notification

import (
	"fmt"
	"github.com/wneessen/go-mail"
	"os"
	"strings"
)

var carrierDomains = map[string]string{
	"att":        "txt.att.net",
	"verizon":    "vtext.com",
	"tmobile":    "tmomail.net",
	"sprint":     "messaging.sprintpcs.com",
	"boost":      "sms.myboostmobile.com",
	"cricket":    "sms.cricketwireless.net",
	"metropcs":   "mymetropcs.com",
	"uscellular": "email.uscc.net",
	"virgin":     "vmobl.com",
}

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

func toEmail(recipient Recipient) (string, error) {
	domain, found := carrierDomains[strings.ToLower(recipient.Carrier)]
	if !found {
		return "", fmt.Errorf("carrier domain not found for %s", recipient.Carrier)
	}
	return recipient.PhoneNumber + "@" + domain, nil
}

func (e *emailNotifier) Text(sms Sms, recipient Recipient) error {
	message := mail.NewMsg()
	if err := message.From(e.smtpUser); err != nil {
		return fmt.Errorf("failed to set FROM address: %w", err)
	}
	recipientEmail, err := toEmail(recipient)
	if err != nil {
		return fmt.Errorf("failed to convert recipient to email: %w", err)
	}
	if err := message.To(recipientEmail); err != nil {
		return fmt.Errorf("failed to set TO address: %w", err)
	}
	message.Subject("")
	message.SetBodyString(mail.TypeTextPlain, sms.Body)

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
