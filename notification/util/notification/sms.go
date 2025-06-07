package util

import (
	"fmt"

	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/notification"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSMS(event common.NotificationEvent, cfg notification.Config) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

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
