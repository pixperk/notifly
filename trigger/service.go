package trigger

import (
	"github.com/nats-io/nats.go"
	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/trigger/util"
)

type Service interface {
	TriggerNotification(event common.NotificationEvent) error
}

type triggerService struct {
	js nats.JetStreamContext
}

func NewService(js nats.JetStreamContext) Service {
	return &triggerService{
		js,
	}
}

func (s *triggerService) TriggerNotification(event common.NotificationEvent) error {
	return util.PublishNotif(s.js, event)
}
