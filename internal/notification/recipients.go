package notification

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

func ParseRecipients(rawRecipient string) []Recipient {
	var response []Recipient
	err := json.Unmarshal([]byte(rawRecipient), &response)
	if err != nil {
		log.Error().Err(err).Msgf("unable to parse recipients %s", rawRecipient)
		return response
	}

	for i, r := range response {
		if r.Carrier != "" && r.PhoneNumber != "" {
			email, err := PhoneToEmail(r)
			if err != nil {
				log.Error().Err(err).Msgf("unable to convert phone number to email %s/%s", r.PhoneNumber, r.Carrier)
				continue
			}
			response[i].Email = email
		}
	}

	return response
}
