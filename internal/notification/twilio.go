package notification

import (
	"fmt"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"os"
)

type (
	twilioNotifier struct {
	}
)

func NewTwilioNotifier() Notifier {
	return &twilioNotifier{}
}

func (t *twilioNotifier) Notify(sms Message, recipient Recipient) error {
	client := twilio.NewRestClient()

	params := &openapi.CreateMessageParams{}
	params.SetTo(recipient.PhoneNumber)
	params.SetFrom(os.Getenv("TWILIO_PHONE_NUMBER"))
	params.SetBody(sms.Body)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}

	return nil
}
