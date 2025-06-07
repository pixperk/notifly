package dispatcher

import (
	"fmt"

	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/notification"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Dispatcher interface {
	Send(event common.NotificationEvent) error
}

type EmailDispatcher struct {
	cfg notification.Config
}

func (d *EmailDispatcher) Send(event common.NotificationEvent) error {
	// Implement email sending logic here
	// For example, use an SMTP client to send the email
	return nil
}

type SMSDispatcher struct {
	cfg notification.Config
}

func (d *SMSDispatcher) Send(event common.NotificationEvent) error {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: d.cfg.TwilioAccountSID,
		Password: d.cfg.TwilioAuthToken,
	})

	params := &openapi.CreateMessageParams{}
	params.SetTo(event.Recipient)
	params.SetFrom(d.cfg.TwilioPhoneNumber)
	params.SetBody(fmt.Sprintf("%s\n%s", event.Subject, event.Body))

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %w", err)
	}
	return nil
}

func GetDispatcher(event common.NotificationEvent, cfg notification.Config) Dispatcher {
	switch event.Type {
	case "EMAIL":
		return &EmailDispatcher{
			cfg,
		}
	case "SMS":
		return &SMSDispatcher{
			cfg,
		}
	default:
		return nil
	}
}
