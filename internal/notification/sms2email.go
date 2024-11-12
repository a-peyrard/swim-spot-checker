package notification

import (
	"fmt"
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

func PhoneToEmail(recipient Recipient) (string, error) {
	domain, found := carrierDomains[strings.ToLower(recipient.Carrier)]
	if !found {
		return "", fmt.Errorf("carrier domain not found for %s", recipient.Carrier)
	}
	return recipient.PhoneNumber + "@" + domain, nil
}
