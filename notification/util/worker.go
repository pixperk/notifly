package util

import (
	"errors"
	"log"

	"github.com/pixperk/notifly/common"
	"github.com/pixperk/notifly/notification"
	"github.com/pixperk/notifly/notification/dispatcher"
)

func processJob(job common.NotificationEvent, cfg notification.Config) error {
	log.Printf("Processing job: %v", job.NotificationId)

	notifDispatcher := dispatcher.GetDispatcher(job, cfg)
	if notifDispatcher == nil {
		return errors.New("no dispatcher found for job type: " + job.Type)
	}

	err := notifDispatcher.Send(job)
	if err != nil {
		return err
	}

	log.Printf("Successfully sent notification for job %v", job.NotificationId)
	return nil
}
