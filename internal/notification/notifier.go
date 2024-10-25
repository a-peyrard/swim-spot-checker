package notification

type (
	Recipient struct {
		Carrier     string
		PhoneNumber string
	}

	Sms struct {
		Body string
	}

	Notifier interface {
		Text(sms Sms, recipient Recipient) error
	}
)
