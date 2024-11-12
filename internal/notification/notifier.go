package notification

type (
	Recipient struct {
		Carrier     string `json:"carrier"`
		PhoneNumber string `json:"phone_number"`
		Email       string `json:"email"`
	}

	Message struct {
		Subject string
		Body    string
	}

	Notifier interface {
		Notify(msg Message, recipient Recipient) error
	}
)
