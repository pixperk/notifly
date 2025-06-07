package dispatcher

import "github.com/pixperk/notifly/common"

type Dispatcher interface {
	Send(event common.NotificationEvent) error
}

type EmailDispatcher struct{}

func (d *EmailDispatcher) Send(event common.NotificationEvent) error {
	// Implement email sending logic here
	// For example, use an SMTP client to send the email
	return nil
}

type SMSDispatcher struct{}

func (d *SMSDispatcher) Send(event common.NotificationEvent) error {
	// Implement SMS sending logic here
	// For example, use an SMS gateway API to send the message
	return nil
}

func GetDispatcher(event common.NotificationEvent) Dispatcher {
	switch event.Type {
	case "EMAIL":
		return &EmailDispatcher{}
	case "SMS":
		return &SMSDispatcher{}
	default:
		return nil
	}
}
