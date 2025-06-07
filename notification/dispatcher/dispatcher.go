package dispatcher

import (
	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/notification"
	util "github.com/pixperk/notifly/notification/util/notification"
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
	return util.SendSMS(event, d.cfg)

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
