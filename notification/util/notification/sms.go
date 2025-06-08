package util

import (
	"fmt"

	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/notification"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

var ErrInvalidPhoneNumber = fmt.Errorf("invalid phone number format")

func SendSMS(event common.NotificationEvent, cfg notification.Config) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

	// Validate phone number
	if !CheckIfPhoneNumberIsValid(event.Recipient) {
		return ErrInvalidPhoneNumber
	}

	params := &openapi.CreateMessageParams{}
	params.SetTo(event.Recipient)
	params.SetFrom(cfg.TwilioPhoneNumber)
	params.SetBody(fmt.Sprintf("%s\n%s", event.Subject, event.Body))

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	return nil
}

func CheckIfPhoneNumberIsValid(phoneNumber string) bool {
	// Validate phone number according to E.164 format guidelines
	// E.164 format: +[country code][area code][local phone number]
	// Example: +1234567890

	// Check if phone number is empty
	if len(phoneNumber) == 0 {
		return false
	}

	// Check if phone number starts with a '+'
	if phoneNumber[0] != '+' {
		return false
	}

	// Check if the remaining characters are digits
	for _, r := range phoneNumber[1:] {
		if r < '0' || r > '9' {
			return false
		}
	}

	// Check if the length is reasonable (typically between 8 and 15 digits plus the '+')
	// Most international numbers are between 9 and 15 digits according to E.164 standard
	if len(phoneNumber) < 9 || len(phoneNumber) > 16 {
		return false
	}

	return true
}
