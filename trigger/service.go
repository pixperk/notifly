package trigger

import (
	"github.com/nats-io/nats.go"
	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/trigger/util"
)

type Service interface {
	TriggerNotification(event common.NotificationEvent) (string, error)
}

type triggerService struct {
	nc *nats.Conn
}

func NewService(nc *nats.Conn) Service {
	return &triggerService{
		nc: nc,
	}
}

func (s *triggerService) TriggerNotification(event common.NotificationEvent) (string, error) {
	return util.PublishNotif(s.nc, event)
}
